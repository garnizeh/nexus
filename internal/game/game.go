package game

import (
	"math"
	"nexussiege/internal/entity"
)

// GameState representa o estado completo do jogo
type GameState struct {
	Tick        uint64
	Gold        float64
	Essence     float64
	Wave        int
	Entities    map[uint64]*entity.Entity
	Towers      map[uint64]*entity.Tower
	Units       map[uint64]*entity.Unit
	SwarmUnits  map[uint64]*entity.SwarmUnit
	Enemies     map[uint64]*entity.Enemy
	Heroes      map[uint64]*entity.Hero
	Projectiles map[uint64]*entity.Projectile
	NextID      uint64
}

// NewGameState cria um novo estado de jogo
func NewGameState() *GameState {
	return &GameState{
		Tick:        0,
		Gold:        100,
		Essence:     50, // Começa com essência para o Enxame
		Wave:        0,
		Entities:    make(map[uint64]*entity.Entity),
		Towers:      make(map[uint64]*entity.Tower),
		Units:       make(map[uint64]*entity.Unit),
		SwarmUnits:  make(map[uint64]*entity.SwarmUnit),
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
func (gs *GameState) AddTower(x, y float64, towerType string) *entity.Tower {
	id := gs.GenerateID()
	tower := &entity.Tower{
		Entity: entity.Entity{
			ID:        id,
			Type:      entity.EntityTower,
			Faction:   entity.FactionIronVanguard,
			X:         x,
			Y:         y,
			Health:    100,
			MaxHealth: 100,
			Dead:      false,
		},
		Level:     1,
		TowerType: towerType,
	}
	
	// Configurar stats baseados no tipo de torre
	switch towerType {
	case "cannon":
		tower.Damage = 25
		tower.Range = 80
		tower.FireRate = 0.5
	case "laser":
		tower.Damage = 5
		tower.Range = 120
		tower.FireRate = 2.0
	case "missile":
		tower.Damage = 15
		tower.Range = 150
		tower.FireRate = 0.8
	default: // basic
		tower.Damage = 10
		tower.Range = 100
		tower.FireRate = 1.0
	}
	
	gs.Towers[id] = tower
	gs.Entities[id] = &tower.Entity
	return tower
}

// AddUnit adiciona uma unidade controlável ao jogo
func (gs *GameState) AddUnit(x, y float64, unitType string, faction entity.Faction) *entity.Unit {
	id := gs.GenerateID()
	unit := &entity.Unit{
		Entity: entity.Entity{
			ID:        id,
			Type:      entity.EntityUnit,
			Faction:   faction,
			X:         x,
			Y:         y,
			Health:    100,
			MaxHealth: 100,
			Dead:      false,
		},
		UnitType:    unitType,
		TargetX:     x,
		TargetY:     y,
		Moving:      false,
		Attacking:   false,
	}
	
	// Configurar stats baseados no tipo de unidade
	switch unitType {
	case "tank":
		unit.Damage = 15
		unit.Speed = 30
		unit.AttackRange = 40
		unit.Health = 200
		unit.MaxHealth = 200
	case "ranger":
		unit.Damage = 12
		unit.Speed = 50
		unit.AttackRange = 120
		unit.Health = 60
		unit.MaxHealth = 60
	default: // soldier
		unit.Damage = 10
		unit.Speed = 40
		unit.AttackRange = 50
	}
	
	gs.Units[id] = unit
	gs.Entities[id] = &unit.Entity
	return unit
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
	delete(gs.SwarmUnits, id)
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
	
	// Atualizar unidades da Vanguarda
	for _, unit := range gs.Units {
		gs.updateUnit(unit, dt)
	}
	
	// Atualizar unidades do Enxame
	for _, swarmUnit := range gs.SwarmUnits {
		gs.updateSwarmUnit(swarmUnit, dt)
	}
	
	// Atualizar heróis
	for _, hero := range gs.Heroes {
		gs.updateHero(hero, dt)
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
	
	for id, unit := range gs.Units {
		if unit.Dead {
			toRemove = append(toRemove, id)
		}
	}
	
	for id, swarmUnit := range gs.SwarmUnits {
		if swarmUnit.Dead {
			toRemove = append(toRemove, id)
		}
	}
	
	for _, id := range toRemove {
		gs.RemoveEntity(id)
	}
}

// updateUnit move e faz a unidade combater
func (gs *GameState) updateUnit(unit *entity.Unit, dt float64) {
	// Se estiver movendo, mover até o alvo
	if unit.Moving {
		dx := unit.TargetX - unit.X
		dy := unit.TargetY - unit.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist < 1 {
			unit.Moving = false
		} else {
			unit.X += dx / dist * unit.Speed * dt
			unit.Y += dy / dist * unit.Speed * dt
		}
	}
	
	// Procurar inimigos próximos para atacar
	unit.Attacking = false
	for _, enemy := range gs.Enemies {
		if enemy.Dead {
			continue
		}
		dx := enemy.X - unit.X
		dy := enemy.Y - unit.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist <= unit.AttackRange {
			// Atacar inimigo
			unit.Attacking = true
			enemy.Health -= unit.Damage * dt
			if enemy.Health <= 0 {
				enemy.Dead = true
				gs.Gold += enemy.Bounty
			}
			break
		}
	}
	
	// Se não atacou inimigo, procurar unidades inimigas
	if !unit.Attacking {
		for _, otherUnit := range gs.Units {
			if otherUnit.Dead || otherUnit.ID == unit.ID {
				continue
			}
			// Unidades de facções diferentes se atacam
			if unit.Faction != otherUnit.Faction && otherUnit.Faction != entity.FactionNone {
				dx := otherUnit.X - unit.X
				dy := otherUnit.Y - unit.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				
				if dist <= unit.AttackRange {
					unit.Attacking = true
					otherUnit.Health -= unit.Damage * dt
					if otherUnit.Health <= 0 {
						otherUnit.Dead = true
					}
					break
				}
			}
		}
	}
}

// updateSwarmUnit move e faz a unidade do Enxame combater
func (gs *GameState) updateSwarmUnit(swarmUnit *entity.SwarmUnit, dt float64) {
	// Se estiver movendo, mover até o alvo
	if swarmUnit.Moving {
		dx := swarmUnit.TargetX - swarmUnit.X
		dy := swarmUnit.TargetY - swarmUnit.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist < 1 {
			swarmUnit.Moving = false
		} else {
			swarmUnit.X += dx / dist * swarmUnit.Speed * dt
			swarmUnit.Y += dy / dist * swarmUnit.Speed * dt
		}
	}
	
	// Procurar inimigos próximos para atacar (torres e unidades da Vanguarda)
	swarmUnit.Attacking = false
	
	// Atacar torres
	for _, tower := range gs.Towers {
		if tower.Dead {
			continue
		}
		dx := tower.X - swarmUnit.X
		dy := tower.Y - swarmUnit.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist <= swarmUnit.AttackRange {
			swarmUnit.Attacking = true
			tower.Health -= swarmUnit.Damage * dt
			if tower.Health <= 0 {
				tower.Dead = true
				gs.Essence += 5 // Ganha essência ao destruir torre
			}
			break
		}
	}
	
	// Se não atacou torre, atacar unidades da Vanguarda
	if !swarmUnit.Attacking {
		for _, otherUnit := range gs.Units {
			if otherUnit.Dead {
				continue
			}
			dx := otherUnit.X - swarmUnit.X
			dy := otherUnit.Y - swarmUnit.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			
			if dist <= swarmUnit.AttackRange {
				swarmUnit.Attacking = true
				otherUnit.Health -= swarmUnit.Damage * dt
				if otherUnit.Health <= 0 {
					otherUnit.Dead = true
					gs.Essence += 3 // Ganha essência ao matar unidade
				}
				break
			}
		}
	}
	
	// Se não atacou nada, atacar heróis
	if !swarmUnit.Attacking {
		for _, hero := range gs.Heroes {
			if hero.Dead {
				continue
			}
			dx := hero.X - swarmUnit.X
			dy := hero.Y - swarmUnit.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			
			if dist <= swarmUnit.AttackRange {
				swarmUnit.Attacking = true
				hero.Health -= swarmUnit.Damage * dt
				if hero.Health <= 0 {
					hero.Dead = true
					gs.Essence += 20 // Ganha muita essência ao matar herói
				}
				break
			}
		}
	}
}

// updateHero atualiza o herói
func (gs *GameState) updateHero(hero *entity.Hero, dt float64) {
	// Atualizar cooldown de habilidades
	if hero.AbilityCooldown > 0 {
		hero.AbilityCooldown -= dt
		if hero.AbilityCooldown <= 0 {
			hero.CanUseAbility = true
			hero.AbilityCooldown = 0
		}
	}
	
	// Herói segue lógica similar às unidades
	// (movimento e combate serão implementados no frontend)
}

// AddSwarmUnit adiciona uma unidade do Enxame ao jogo
func (gs *GameState) AddSwarmUnit(x, y float64, unitType entity.SwarmUnitType) *entity.SwarmUnit {
	id := gs.GenerateID()
	swarmUnit := entity.NewSwarmUnit(id, x, y, unitType)
	
	gs.SwarmUnits[id] = swarmUnit
	gs.Entities[id] = &swarmUnit.Entity
	return swarmUnit
}

// AddHero adiciona um herói ao jogo
func (gs *GameState) AddHero(x, y float64, heroClass entity.HeroClass) *entity.Hero {
	id := gs.GenerateID()
	hero := entity.NewHero(id, x, y, heroClass)
	
	gs.Heroes[id] = hero
	gs.Entities[id] = &hero.Entity
	return hero
}

// UseHeroAbility usa uma habilidade de um herói
func (gs *GameState) UseHeroAbility(heroID uint64, ability entity.AbilityType, targetX, targetY float64) bool {
	hero, exists := gs.Heroes[heroID]
	if !exists || hero.Dead {
		return false
	}
	
	if !hero.UseAbility(ability) {
		return false
	}
	
	// Aplicar efeito da habilidade
	switch ability {
	case entity.AbilityBlast:
		// Dano em área ao redor do herói
		for _, enemy := range gs.Enemies {
			dx := enemy.X - hero.X
			dy := enemy.Y - hero.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= 100 { // Raio de 100 unidades
				enemy.Health -= 50
				if enemy.Health <= 0 {
					enemy.Dead = true
					gs.Gold += enemy.Bounty
					hero.GainExperience(10)
				}
			}
		}
		// Também afeta unidades do Enxame
		for _, swarmUnit := range gs.SwarmUnits {
			dx := swarmUnit.X - hero.X
			dy := swarmUnit.Y - hero.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= 100 {
				swarmUnit.Health -= 50
				if swarmUnit.Health <= 0 {
					swarmUnit.Dead = true
					hero.GainExperience(15)
				}
			}
		}
		
	case entity.AbilityHeal:
		// Cura todas as unidades aliadas próximas
		for _, unit := range gs.Units {
			if unit.Faction == entity.FactionIronVanguard && !unit.Dead {
				dx := unit.X - hero.X
				dy := unit.Y - hero.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist <= 80 {
					unit.Health = unit.MaxHealth
				}
			}
		}
		// Cura o próprio herói
		hero.Health = hero.MaxHealth
		
	case entity.AbilityShield:
		// Escudo temporário (implementação simplificada - cura imediata)
		hero.Health = hero.MaxHealth * 1.5
		
	case entity.AbilityTeleport:
		// Teleporta para a posição alvo
		hero.X = targetX
		hero.Y = targetY
		
	case entity.AbilityNuke:
		// Dano massivo em área grande
		for _, enemy := range gs.Enemies {
			dx := enemy.X - targetX
			dy := enemy.Y - targetY
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= 150 {
				enemy.Health -= 100
				if enemy.Health <= 0 {
					enemy.Dead = true
					gs.Gold += enemy.Bounty
					hero.GainExperience(20)
				}
			}
		}
		// Também afeta unidades do Enxame
		for _, swarmUnit := range gs.SwarmUnits {
			dx := swarmUnit.X - targetX
			dy := swarmUnit.Y - targetY
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= 150 {
				swarmUnit.Health -= 100
				if swarmUnit.Health <= 0 {
					swarmUnit.Dead = true
					hero.GainExperience(25)
				}
			}
		}
	}
	
	return true
}

// SabotageWave envia uma horda de sabotagem contra o jogador
func (gs *GameState) SabotageWave(waveIntensity int) int {
	enemiesSpawned := 0
	
	// Spawn de inimigos baseado na intensidade
	for i := 0; i < waveIntensity; i++ {
		x := float64(i % 10 * 50)
		y := -float64(i/10 * 50)
		gs.AddEnemy(x, y, gs.Wave+1)
		enemiesSpawned++
	}
	
	// Chance de spawnar unidade do Enxame
	if waveIntensity >= 5 {
		swarmCount := waveIntensity / 5
		for i := 0; i < swarmCount; i++ {
			x := float64(400 + i*30)
			y := -50.0
			var unitType entity.SwarmUnitType
			if i%3 == 0 {
				unitType = entity.SwarmBehemoth
			} else if i%3 == 1 {
				unitType = entity.SwarmStalker
			} else {
				unitType = entity.SwarmDrone
			}
			gs.AddSwarmUnit(x, y, unitType)
		}
	}
	
	gs.Wave++
	return enemiesSpawned
}
