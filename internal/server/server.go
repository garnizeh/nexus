package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"nexussiege/internal/entity"
	"nexussiege/internal/game"
)

// Server representa o servidor do jogo
type Server struct {
	gameState *game.GameState
	mu        sync.RWMutex
	httpServer *http.Server
}

// NewServer cria um novo servidor
func NewServer() *Server {
	return &Server{
		gameState: game.NewGameState(),
	}
}

// Start inicia o servidor HTTP e o game loop
func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	
	// Servir arquivos estáticos
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	
	// API endpoints
	mux.HandleFunc("/api/state", s.handleGetState)
	mux.HandleFunc("/api/build", s.handleBuildTower)
	mux.HandleFunc("/api/wave", s.handleStartWave)
	mux.HandleFunc("/api/unit", s.handleCreateUnit)
	mux.HandleFunc("/api/move", s.handleMoveUnit)
	mux.HandleFunc("/api/sabotage", s.handleSabotage)
	
	// Página principal
	mux.HandleFunc("/", s.handleIndex)
	
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	
	// Iniciar game loop em goroutine
	go s.gameLoop()
	
	log.Printf("Servidor iniciando em %s", addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown para o servidor graciosamente
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// gameLoop roda a 30Hz (tick rate do servidor)
func (s *Server) gameLoop() {
	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()
	
	for range ticker.C {
		s.mu.Lock()
		s.gameState.Update(1.0 / 30.0)
		s.mu.Unlock()
	}
}

// handleIndex serve a página principal
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}

// handleGetState retorna o estado atual do jogo
func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.gameState)
}

// handleBuildTower lida com a construção de torres
func (s *Server) handleBuildTower(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		X         float64 `json:"x"`
		Y         float64 `json:"y"`
		TowerType string  `json:"tower_type"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Verificar se tem ouro suficiente
	cost := 50.0
	switch req.TowerType {
	case "cannon":
		cost = 80
	case "laser":
		cost = 70
	case "missile":
		cost = 100
	}
	
	if s.gameState.Gold < cost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Gold insufficient"})
		return
	}
	
	s.gameState.Gold -= cost
	tower := s.gameState.AddTower(req.X, req.Y, req.TowerType)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"tower_id": tower.ID,
		"cost":     cost,
	})
}

// handleStartWave inicia uma nova wave de inimigos
func (s *Server) handleStartWave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.gameState.Wave++
	wave := s.gameState.Wave
	
	// Spawnar inimigos baseado na wave
	numEnemies := 5 + wave*2
	for i := 0; i < numEnemies; i++ {
		// Spawn em posição inicial do caminho
		s.gameState.AddEnemy(0, float64(i*20), wave)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"wave": wave,
		"enemies": numEnemies,
	})
}

// handleCreateUnit lida com a criação de unidades
func (s *Server) handleCreateUnit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
		UnitType string  `json:"unit_type"` // soldier, tank, ranger
		Faction  int     `json:"faction"`   // 1 = IronVanguard, 2 = VoraciousSwarm
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Verificar se tem ouro suficiente
	cost := 60.0
	switch req.UnitType {
	case "tank":
		cost = 100
	case "ranger":
		cost = 70
	}
	
	if s.gameState.Gold < cost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Gold insufficient"})
		return
	}
	
	s.gameState.Gold -= cost
	faction := entity.Faction(req.Faction)
	if faction == entity.FactionNone {
		faction = entity.FactionIronVanguard
	}
	unit := s.gameState.AddUnit(req.X, req.Y, req.UnitType, faction)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"unit_id": unit.ID,
		"cost":    cost,
	})
}

// handleMoveUnit lida com o movimento de unidades
func (s *Server) handleMoveUnit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		UnitID  uint64  `json:"unit_id"`
		TargetX float64 `json:"target_x"`
		TargetY float64 `json:"target_y"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	unit, exists := s.gameState.Units[req.UnitID]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unit not found"})
		return
	}
	
	unit.TargetX = req.TargetX
	unit.TargetY = req.TargetY
	unit.Moving = true
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// handleSabotage lida com o envio de horda de sabotagem
func (s *Server) handleSabotage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		EssenceCost float64 `json:"essence_cost"`
		TargetX     float64 `json:"target_x"`
		TargetY     float64 `json:"target_y"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Verificar se tem essência suficiente
	if s.gameState.Essence < req.EssenceCost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Essence insufficient"})
		return
	}
	
	s.gameState.Essence -= req.EssenceCost
	
	// Spawnar horda de inimigos na posição alvo (sabotagem)
	numEnemies := int(req.EssenceCost / 10)
	for i := 0; i < numEnemies; i++ {
		s.gameState.AddEnemy(req.TargetX+float64(i*5), req.TargetY, 99)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"enemies":  numEnemies,
		"message":  "Horda de sabotagem enviada!",
	})
}
