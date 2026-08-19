package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

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
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Verificar se tem ouro suficiente
	if s.gameState.Gold < 50 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Gold insufficient"})
		return
	}
	
	s.gameState.Gold -= 50
	tower := s.gameState.AddTower(req.X, req.Y)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"tower_id": tower.ID,
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
