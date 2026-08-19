package entity

// AbilityType define o tipo de habilidade
type AbilityType int

const (
	AbilityBlast AbilityType = iota
	AbilityHeal
	AbilityShield
	AbilityTeleport
	AbilityNuke
)

// HeroClass define a classe do herói
type HeroClass int

const (
	HeroCommander HeroClass = iota
	HeroEngineer
	HeroScout
	HeroJuggernaut
)

// Hero representa um herói com habilidades
type Hero struct {
	Entity
	Damage          float64
	Speed           float64
	Level           int
	Experience      float64
	AbilityCooldown float64
	CanUseAbility   bool
	HeroClass       HeroClass
	CurrentAbility  AbilityType
}

// NewHero cria um novo herói
func NewHero(id uint64, x, y float64, heroClass HeroClass) *Hero {
	hero := &Hero{
		Entity: Entity{
			ID:        id,
			Type:      EntityHero,
			Faction:   FactionIronVanguard,
			X:         x,
			Y:         y,
			Health:    150,
			MaxHealth: 150,
			Dead:      false,
		},
		HeroClass:     heroClass,
		Level:         1,
		Experience:    0,
		Speed:         45,
		Damage:        20,
		CanUseAbility: true,
	}

	// Configurar stats baseados na classe
	switch heroClass {
	case HeroCommander:
		hero.Health = 180
		hero.MaxHealth = 180
		hero.Damage = 25
	case HeroEngineer:
		hero.Health = 120
		hero.MaxHealth = 120
		hero.Damage = 15
		hero.Speed = 40
	case HeroScout:
		hero.Health = 100
		hero.MaxHealth = 100
		hero.Damage = 30
		hero.Speed = 60
	case HeroJuggernaut:
		hero.Health = 250
		hero.MaxHealth = 250
		hero.Damage = 35
		hero.Speed = 30
	}

	return hero
}

// GetAbilities retorna as habilidades disponíveis para esta classe
func (h *Hero) GetAbilities() []AbilityType {
	switch h.HeroClass {
	case HeroCommander:
		return []AbilityType{AbilityBlast, AbilityShield, AbilityHeal}
	case HeroEngineer:
		return []AbilityType{AbilityShield, AbilityHeal, AbilityTeleport}
	case HeroScout:
		return []AbilityType{AbilityBlast, AbilityTeleport, AbilityNuke}
	case HeroJuggernaut:
		return []AbilityType{AbilityBlast, AbilityShield, AbilityNuke}
	default:
		return []AbilityType{AbilityBlast}
	}
}

// UseAbility usa uma habilidade do herói
func (h *Hero) UseAbility(ability AbilityType) bool {
	if !h.CanUseAbility {
		return false
	}

	// Verificar se a habilidade é válida para esta classe
	validAbilities := h.GetAbilities()
	valid := false
	for _, a := range validAbilities {
		if a == ability {
			valid = true
			break
		}
	}

	if !valid {
		return false
	}

	h.CurrentAbility = ability
	h.CanUseAbility = false
	h.AbilityCooldown = 10.0 // 10 segundos de cooldown

	return true
}

// Update atualiza o cooldown das habilidades
func (h *Hero) Update(dt float64) {
	if h.AbilityCooldown > 0 {
		h.AbilityCooldown -= dt
		if h.AbilityCooldown <= 0 {
			h.CanUseAbility = true
			h.AbilityCooldown = 0
		}
	}
}

// GainExperience concede experiência ao herói
func (h *Hero) GainExperience(amount float64) {
	h.Experience += amount

	// Level up a cada 100 XP
	levelUpThreshold := float64(h.Level * 100)
	for h.Experience >= levelUpThreshold {
		h.Level++
		h.Experience -= levelUpThreshold
		h.Health = h.MaxHealth // Cura ao subir de nível
		h.MaxHealth *= 1.1     // Aumenta vida máxima
		h.Damage *= 1.1        // Aumenta dano
		levelUpThreshold = float64(h.Level * 100)
	}
}
