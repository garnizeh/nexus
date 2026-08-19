package game

import (
	"math"
	"nexussiege/internal/entity"
)

// GameState representa o estado completo do jogo
type GameState struct {
	Tick        uint64
	Gold        float64
	Essence     float64 // recurso da facção Enxame
	Wave        int
	Entities    map[uint64]*entity.Entity
	Towers      map[uint64]*entity.Tower
	Units       map[uint64]*entity.Unit
	Enemies     map[uint64]*entity.Enemy
	Heroes      map[uint64]*entity.Hero
	SwarmUnits  map[uint64]*entity.SwarmUnit
	Projectiles map[uint64]*entity.Projectile
	NextID      uint64
}

// NewGameState cria um novo estado de jogo
func NewGameState() *GameState {
	return &GameState{
		Tick:        0,
		Gold:        100,
		Essence:     50, // começa com essência para o Enxame
		Wave:        0,
		Entities:    make(map[uint64]*entity.Entity),
		Towers:      make(map[uint64]*entity.Tower),
		Units:       make(map[uint64]*entity.Unit),
		Enemies:     make(map[uint64]*entity.Enemy),
		Heroes:      make(map[uint64]*entity.Hero),
		SwarmUnits:  make(map[uint64]*entity.SwarmUnit),
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

// AddHero adiciona um herói ao jogo
func (gs *GameState) AddHero(x, y float64, heroType string, abilityType entity.AbilityType) *entity.Hero {
	id := gs.GenerateID()
	hero := &entity.Hero{
		Entity: entity.Entity{
			ID:        id,
			Type:      entity.EntityHero,
			X:         x,
			Y:         y,
			Health:    150,
			MaxHealth: 150,
			Dead:      false,
		},
		HeroType: heroType,
		Damage:   20,
		Speed:    45,
		Level:    1,
		Experience: 0,
		Ability: entity.HeroAbility{
			Type:            abilityType,
			CurrentCooldown: 0,
			Ready:           true,
			Level:           1,
		},
		CanUseAbility: true,
	}
	
	// Configurar stats baseados no tipo de herói
	switch heroType {
	case "tank":
		hero.Health = 250
		hero.MaxHealth = 250
		hero.Damage = 15
	case "support":
		hero.Health = 100
		hero.MaxHealth = 100
		hero.Damage = 10
	case "damage":
		hero.Health = 120
		hero.MaxHealth = 120
		hero.Damage = 30
	case "swarm":
		hero.Faction = entity.FactionVoraciousSwarm
		hero.Health = 180
		hero.MaxHealth = 180
		hero.Damage = 25
	}
	
	gs.Heroes[id] = hero
	gs.Entities[id] = &hero.Entity
	return hero
}

// AddSwarmUnit adiciona uma unidade do Enxame ao jogo
func (gs *GameState) AddSwarmUnit(x, y float64, unitType string) *entity.SwarmUnit {
	id := gs.GenerateID()
	swarm := &entity.SwarmUnit{
		Entity: entity.Entity{
			ID:        id,
			Type:      entity.EntityUnit,
			Faction:   entity.FactionVoraciousSwarm,
			X:         x,
			Y:         y,
			Health:    40,
			MaxHealth: 40,
			Dead:      false,
		},
		UnitType:    unitType,
		TargetX:     x,
		TargetY:     y,
		Moving:      false,
		Attacking:   false,
	}
	
	// Configurar stats baseados no tipo de unidade do Enxame
	switch unitType {
	case "behemoth":
		swarm.Damage = 20
		swarm.Speed = 25
		swarm.AttackRange = 30
		swarm.Health = 150
		swarm.MaxHealth = 150
		swarm.EssenceCost = 40
	case "stalker":
		swarm.Damage = 12
		swarm.Speed = 55
		swarm.AttackRange = 40
		swarm.Health = 50
		swarm.MaxHealth = 50
		swarm.EssenceCost = 25
	default: // drone
		swarm.Damage = 8
		swarm.Speed = 45
		swarm.AttackRange = 35
		swarm.EssenceCost = 15
	}
	
	gs.SwarmUnits[id] = swarm
	gs.Entities[id] = &swarm.Entity
	return swarm
}

// UseHeroAbility usa a habilidade do herói
func (gs *GameState) UseHeroAbility(heroID uint64) bool {
	hero, exists := gs.Heroes[heroID]
	if !exists || !hero.CanUseAbility {
		return false
	}
	
	ability := entity.GetAbility(hero.Ability.Type)
	if ability.Cooldown > 0 {
		hero.Ability.CurrentCooldown = ability.Cooldown
		hero.Ability.Ready = false
		hero.CanUseAbility = false
	}
	
	// Aplicar efeito da habilidade baseado no tipo
	switch hero.Ability.Type {
	case entity.AbilityShockwave:
		// Dano em área ao redor do herói
		for _, enemy := range gs.Enemies {
			dx := enemy.X - hero.X
			dy := enemy.Y - hero.Y
			dist := dx*dx + dy*dy
			if dist <= ability.Range*ability.Range {
				enemy.Health -= ability.Damage
				if enemy.Health <= 0 {
					enemy.Dead = true
					gs.Gold += enemy.Bounty
				}
			}
		}
		// Afeta unidades do Enxame também
		for _, swarm := range gs.SwarmUnits {
			dx := swarm.X - hero.X
			dy := swarm.Y - hero.Y
			dist := dx*dx + dy*dy
			if dist <= ability.Range*ability.Range {
				swarm.Health -= ability.Damage
				if swarm.Health <= 0 {
					swarm.Dead = true
					gs.Gold += 5
				}
			}
		}
		
	case entity.AbilityHeal:
		// Cura unidades aliadas próximas
		for _, unit := range gs.Units {
			if unit.Faction == hero.Faction || hero.Faction == entity.FactionNone {
				dx := unit.X - hero.X
				dy := unit.Y - hero.Y
				dist := dx*dx + dy*dy
				if dist <= ability.Range*ability.Range {
					unit.Health += ability.HealAmount
					if unit.Health > unit.MaxHealth {
						unit.Health = unit.MaxHealth
					}
				}
			}
		}
		// Cura o próprio herói
		hero.Health += ability.HealAmount
		if hero.Health > hero.MaxHealth {
			hero.Health = hero.MaxHealth
		}
		
	case entity.AbilitySnipe:
		// Dano massivo no inimigo mais próximo
		var closest *entity.Enemy
		closestDist := ability.Range * ability.Range
		for _, enemy := range gs.Enemies {
			dx := enemy.X - hero.X
			dy := enemy.Y - hero.Y
			dist := dx*dx + dy*dy
			if dist < closestDist && !enemy.Dead {
				closest = enemy
				closestDist = dist
			}
		}
		if closest != nil {
			closest.Health -= ability.Damage
			if closest.Health <= 0 {
				closest.Dead = true
				gs.Gold += closest.Bounty
			}
		}
		
	case entity.AbilitySwarmRush:
		// Aumenta velocidade das unidades do Enxame próximas
		if hero.Faction == entity.FactionVoraciousSwarm {
			for _, swarm := range gs.SwarmUnits {
				dx := swarm.X - hero.X
				dy := swarm.Y - hero.Y
				dist := dx*dx + dy*dy
				if dist <= ability.Range*ability.Range {
					swarm.Speed *= 1.5 // +50% velocidade
				}
			}
		}
		
	case entity.AbilityAcidSpray:
		// Dano contínuo em área
		for _, enemy := range gs.Enemies {
			dx := enemy.X - hero.X
			dy := enemy.Y - hero.Y
			dist := dx*dx + dy*dy
			if dist <= ability.Range*ability.Range {
				enemy.Health -= ability.Damage
				if enemy.Health <= 0 {
					enemy.Dead = true
					gs.Gold += enemy.Bounty
				}
			}
		}
	}
	
	return true
}

// RemoveEntity remove uma entidade do jogo
func (gs *GameState) RemoveEntity(id uint64) {
	delete(gs.Entities, id)
	delete(gs.Towers, id)
	delete(gs.Units, id)
	delete(gs.Enemies, id)
	delete(gs.Heroes, id)
	delete(gs.SwarmUnits, id)
	delete(gs.Projectiles, id)
}

// Update atualiza o estado do jogo
func (gs *GameState) Update(dt float64) {
	gs.Tick++
	
	// Atualizar torres (combate)
	for _, tower := range gs.Towers {
		gs.updateTower(tower, dt)
	}
	
	// Atualizar unidades
	for _, unit := range gs.Units {
		gs.updateUnit(unit, dt)
	}
	
	// Atualizar unidades do Enxame
	for _, swarm := range gs.SwarmUnits {
		gs.updateSwarmUnit(swarm, dt)
	}
	
	// Atualizar heróis
	for _, hero := range gs.Heroes {
		gs.updateHero(hero, dt)
	}
	
	// Atualizar projéteis
	for _, proj := range gs.Projectiles {
		gs.updateProjectile(proj, dt)
	}
	
	// Gerar essência passivamente para o Enxame
	gs.Essence += 0.1 * dt
	
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
	
	for id, swarm := range gs.SwarmUnits {
		if swarm.Dead {
			toRemove = append(toRemove, id)
		}
	}
	
	for id, hero := range gs.Heroes {
		if hero.Dead {
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
	
	// Atacar unidades do Enxame se estiverem próximas
	if !unit.Attacking {
		for _, swarm := range gs.SwarmUnits {
			if swarm.Dead {
				continue
			}
			dx := swarm.X - unit.X
			dy := swarm.Y - unit.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			
			if dist <= unit.AttackRange {
				unit.Attacking = true
				swarm.Health -= unit.Damage * dt
				if swarm.Health <= 0 {
					swarm.Dead = true
					gs.Gold += 5 // recompensa por matar unidade do enxame
				}
				break
			}
		}
	}
}

// updateSwarmUnit move e faz a unidade do Enxame combater
func (gs *GameState) updateSwarmUnit(swarm *entity.SwarmUnit, dt float64) {
	// Se estiver movendo, mover até o alvo
	if swarm.Moving {
		dx := swarm.TargetX - swarm.X
		dy := swarm.TargetY - swarm.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist < 1 {
			swarm.Moving = false
		} else {
			swarm.X += dx / dist * swarm.Speed * dt
			swarm.Y += dy / dist * swarm.Speed * dt
		}
	}
	
	// Procurar inimigos próximos para atacar (torres, unidades, heróis)
	swarm.Attacking = false
	
	// Atacar torres
	for _, tower := range gs.Towers {
		if tower.Dead {
			continue
		}
		dx := tower.X - swarm.X
		dy := tower.Y - swarm.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		
		if dist <= swarm.AttackRange {
			swarm.Attacking = true
			tower.Health -= swarm.Damage * dt
			if tower.Health <= 0 {
				tower.Dead = true
				gs.Essence += 10 // ganha essência ao destruir torre
			}
			break
		}
	}
	
	// Atacar unidades inimigas
	if !swarm.Attacking {
		for _, unit := range gs.Units {
			if unit.Dead {
				continue
			}
			dx := unit.X - swarm.X
			dy := unit.Y - swarm.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			
			if dist <= swarm.AttackRange {
				swarm.Attacking = true
				unit.Health -= swarm.Damage * dt
				if unit.Health <= 0 {
					unit.Dead = true
					gs.Essence += 8
				}
				break
			}
		}
	}
	
	// Atacar heróis
	if !swarm.Attacking {
		for _, hero := range gs.Heroes {
			if hero.Dead {
				continue
			}
			dx := hero.X - swarm.X
			dy := hero.Y - swarm.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			
			if dist <= swarm.AttackRange {
				swarm.Attacking = true
				hero.Health -= swarm.Damage * dt
				if hero.Health <= 0 {
					hero.Dead = true
					gs.Essence += 25
				}
				break
			}
		}
	}
}

// updateHero atualiza o herói e suas habilidades
func (gs *GameState) updateHero(hero *entity.Hero, dt float64) {
	// Reduzir cooldown da habilidade
	if hero.Ability.CurrentCooldown > 0 {
		hero.Ability.CurrentCooldown -= dt
		if hero.Ability.CurrentCooldown <= 0 {
			hero.Ability.Ready = true
			hero.Ability.CurrentCooldown = 0
		}
	}
	
	// Herói pode usar habilidade se estiver pronta
	hero.CanUseAbility = hero.Ability.Ready
}
