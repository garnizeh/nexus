package entity

// SwarmUnitType define o tipo de unidade do Enxame
type SwarmUnitType int

const (
	SwarmDrone SwarmUnitType = iota
	SwarmStalker
	SwarmBehemoth
)

// SwarmUnit representa uma unidade do Enxame Voraz
type SwarmUnit struct {
	Entity
	Damage      float64
	Speed       float64
	AttackRange float64
	TargetX     float64
	TargetY     float64
	Moving      bool
	Attacking   bool
	UnitType    SwarmUnitType
	EssenceCost float64
}

// NewSwarmUnit cria uma nova unidade do Enxame
func NewSwarmUnit(id uint64, x, y float64, unitType SwarmUnitType) *SwarmUnit {
	unit := &SwarmUnit{
		Entity: Entity{
			ID:        id,
			Type:      EntityUnit,
			Faction:   FactionVoraciousSwarm,
			X:         x,
			Y:         y,
			Dead:      false,
		},
		UnitType:    unitType,
		TargetX:     x,
		TargetY:     y,
		Moving:      false,
		Attacking:   false,
	}

	// Configurar stats baseados no tipo
	switch unitType {
	case SwarmDrone:
		unit.Damage = 8
		unit.Speed = 60
		unit.AttackRange = 30
		unit.Health = 40
		unit.MaxHealth = 40
		unit.EssenceCost = 10
	case SwarmStalker:
		unit.Damage = 15
		unit.Speed = 45
		unit.AttackRange = 50
		unit.Health = 80
		unit.MaxHealth = 80
		unit.EssenceCost = 25
	case SwarmBehemoth:
		unit.Damage = 30
		unit.Speed = 20
		unit.AttackRange = 40
		unit.Health = 300
		unit.MaxHealth = 300
		unit.EssenceCost = 60
	}

	return unit
}
