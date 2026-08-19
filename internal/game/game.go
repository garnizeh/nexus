package game

import (
	"nexussiege/internal/entity"
)

// GameState representa o estado completo do jogo
type GameState struct {
	Tick      uint64
	Gold      float64
	Essence   float64
	Wave      int
	Entities  map[uint64]*entity.Entity
	Towers    map[uint64]*entity.Tower
	Units     map[uint64]*entity.Unit
	Enemies   map[uint64]*entity.Enemy
	Heroes    map[uint64]*entity.Hero
	Projectiles map[uint64]*entity.Projectile
	NextID    uint64
}

// NewGameState cria um novo estado de jogo
func NewGameState() *GameState {
	return &GameState{
		Tick:        0,
		Gold:        100,
		Essence:     0,
		Wave:        0,
		Entities:    make(map[uint64]*entity.Entity),
		Towers:      make(map[uint64]*entity.Tower),
		Units:       make(map[uint64]*entity.Unit),
		Enemies:     make(map[uint64]*entity.Enemy),
		Heroes:      make(map[uint64]*entity.Hero),
		Projectiles: make(map[uint64]*entity.Projectile),
		NextID:      1,
	}
}

// GenerateID gera um novo ID único para entidades
func (gs *GameState) GenerateID() uint64 {
	id := gs.NextID
	gs.NextID++
	return id
}

// AddTower adiciona uma torre ao jogo
func (gs *GameState) AddTower(x, y float64) *entity.Tower {
	id := gs.GenerateID()
	tower := &entity.Tower{
		Entity: entity.Entity{
			ID:        id,
			Type:      entity.EntityTower,
			X:         x,
			Y:         y,
			Health:    100,
			MaxHealth: 100,
			Dead:      false,
		},
		Level:    1,
		Damage:   10,
		Range:    100,
		FireRate: 1.0,
		LastFire: 0,
	}
	gs.Towers[id] = tower
	gs.Entities[id] = &tower.Entity
	return tower
}

// AddEnemy adiciona um inimigo ao jogo
func (gs *GameState) AddEnemy(x, y float64, wave int) *entity.Enemy {
	id := gs.GenerateID()
	enemy := &entity.Enemy{
		Entity: entity.Entity{
			ID:        id,
			Type:      entity.EntityEnemy,
			X:         x,
			Y:         y,
			Health:    50,
			MaxHealth: 50,
			Dead:      false,
		},
		Speed:     20,
		Damage:    5,
		Bounty:    10,
		Wave:      wave,
		PathIndex: 0,
	}
	gs.Enemies[id] = enemy
	gs.Entities[id] = &enemy.Entity
	return enemy
}

// RemoveEntity remove uma entidade do jogo
func (gs *GameState) RemoveEntity(id uint64) {
	delete(gs.Entities, id)
	delete(gs.Towers, id)
	delete(gs.Units, id)
	delete(gs.Enemies, id)
	delete(gs.Heroes, id)
	delete(gs.Projectiles, id)
}

// Update atualiza o estado do jogo
func (gs *GameState) Update(dt float64) {
	gs.Tick++
	
	// Atualizar torres (combate)
	for _, tower := range gs.Towers {
		gs.updateTower(tower, dt)
	}
	
	// Atualizar projéteis
	for _, proj := range gs.Projectiles {
		gs.updateProjectile(proj, dt)
	}
	
	// Limpar entidades mortas
	gs.cleanupDeadEntities()
}

// updateTower faz a torre atirar em inimigos próximos
func (gs *GameState) updateTower(tower *entity.Tower, dt float64) {
	tower.LastFire += dt
	
	// Verificar se pode atirar
	if tower.LastFire < 1.0/tower.FireRate {
		return
	}
	
	// Procurar inimigo no alcance
	for _, enemy := range gs.Enemies {
		dx := enemy.X - tower.X
		dy := enemy.Y - tower.Y
		dist := dx*dx + dy*dy
		
		if dist <= tower.Range*tower.Range {
			// Atirar!
			gs.fireProjectile(tower, enemy)
			tower.LastFire = 0
			break
		}
	}
}

// fireProjectile cria um projétil da torre para o inimigo
func (gs *GameState) fireProjectile(tower *entity.Tower, target *entity.Enemy) {
	id := gs.GenerateID()
	proj := &entity.Projectile{
		Entity: entity.Entity{
			ID:   id,
			Type: entity.EntityProjectile,
			X:    tower.X,
			Y:    tower.Y,
		},
		Damage:   tower.Damage,
		Speed:    300,
		TargetID: target.ID,
		SourceID: tower.ID,
	}
	gs.Projectiles[id] = proj
	gs.Entities[id] = &proj.Entity
}

// updateProjectile move o projétil até o alvo
func (gs *GameState) updateProjectile(proj *entity.Projectile, dt float64) {
	target, exists := gs.Enemies[proj.TargetID]
	if !exists || target.Dead {
		proj.Dead = true
		return
	}
	
	dx := target.X - proj.X
	dy := target.Y - proj.Y
	dist := dx*dx + dy*dy
	
	if dist < 1 {
		// Acertou!
		target.Health -= proj.Damage
		if target.Health <= 0 {
			target.Dead = true
			// Conceder ouro
			gs.Gold += target.Bounty
		}
		proj.Dead = true
		return
	}
	
	// Mover projétil
	normalizedDist := dist / (proj.Speed * proj.Speed)
	proj.X += dx / normalizedDist * dt
	proj.Y += dy / normalizedDist * dt
}

// cleanupDeadEntities remove entidades mortas
func (gs *GameState) cleanupDeadEntities() {
	var toRemove []uint64
	
	for id, proj := range gs.Projectiles {
		if proj.Dead {
			toRemove = append(toRemove, id)
		}
	}
	
	for id, enemy := range gs.Enemies {
		if enemy.Dead {
			toRemove = append(toRemove, id)
		}
	}
	
	for _, id := range toRemove {
		gs.RemoveEntity(id)
	}
}
