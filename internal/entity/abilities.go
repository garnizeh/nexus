package entity

// AbilityType define o tipo de habilidade
type AbilityType int

const (
	AbilityNone AbilityType = iota
	AbilityShockwave     // Herói Tanque - onda de choque
	AbilityHeal          // Herói Suporte - cura área
	AbilitySnipe         // Herói Dano - tiro preciso
	AbilitySwarmRush     // Enxame - investida em massa
	AbilityAcidSpray     // Enxame - spray de ácido
)

// Ability representa uma habilidade de herói ou unidade especial
type Ability struct {
	Type        AbilityType
	Name        string
	Description string
	Cooldown    float64 // segundos
	Duration    float64 // duração do efeito
	Range       float64 // alcance da habilidade
	Damage      float64 // dano da habilidade
	HealAmount  float64 // quantidade de cura
	Cost        float64 // custo de recurso (ouro/essência)
}

// GetAbility retorna os dados de uma habilidade pelo tipo
func GetAbility(t AbilityType) Ability {
	switch t {
	case AbilityShockwave:
		return Ability{
			Type:        AbilityShockwave,
			Name:        "Onda de Choque",
			Description: "Derruba todos os inimigos próximos",
			Cooldown:    8.0,
			Duration:    0.5,
			Range:       80.0,
			Damage:      25.0,
			Cost:        0,
		}
	case AbilityHeal:
		return Ability{
			Type:        AbilityHeal,
			Name:        "Campo de Cura",
			Description: "Cura unidades aliadas na área",
			Cooldown:    12.0,
			Duration:    3.0,
			Range:       100.0,
			HealAmount:  30.0,
			Cost:        0,
		}
	case AbilitySnipe:
		return Ability{
			Type:        AbilitySnipe,
			Name:        "Tiro Preciso",
			Description: "Dano massivo em um único alvo",
			Cooldown:    6.0,
			Duration:    0,
			Range:       300.0,
			Damage:      100.0,
			Cost:        0,
		}
	case AbilitySwarmRush:
		return Ability{
			Type:        AbilitySwarmRush,
			Name:        "Investida Voraz",
			Description: "Unidades do enxame ganham velocidade temporária",
			Cooldown:    10.0,
			Duration:    5.0,
			Range:       150.0,
			Cost:        20.0, // custa essência
		}
	case AbilityAcidSpray:
		return Ability{
			Type:        AbilityAcidSpray,
			Name:        "Spray de Ácido",
			Description: "Corrói armaduras e causa dano contínuo",
			Cooldown:    9.0,
			Duration:    4.0,
			Range:       120.0,
			Damage:      15.0, // dano por tick
			Cost:        15.0,
		}
	default:
		return Ability{}
	}
}

// HeroAbility representa uma habilidade equipada em um herói
type HeroAbility struct {
	Type           AbilityType
	CurrentCooldown float64
	Ready          bool
	Level          int
}
