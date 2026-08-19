package entity

// Faction representa a facção da entidade
type Faction int

const (
	FactionNone Faction = iota
	FactionIronVanguard  // Vanguarda de Ferro
	FactionVoraciousSwarm // Enxame Voraz
)

// EntityType define o tipo da entidade
type EntityType int

const (
	EntityTower EntityType = iota
	EntityUnit
	EntityEnemy
	EntityHero
	EntityProjectile
)

// Entity representa qualquer entidade no jogo
type Entity struct {
	ID        uint64
	Type      EntityType
	Faction   Faction
	X         float64
	Y         float64
	Health    float64
	MaxHealth float64
	Dead      bool
}

// Tower representa uma torre
type Tower struct {
	Entity
	Level      int
	Damage     float64
	Range      float64
	FireRate   float64 // tiros por segundo
	LastFire   float64 // tempo do último tiro
	TargetID   uint64
	TowerType  string  // "basic", "cannon", "laser", "missile"
}

// Unit representa uma unidade controlável
type Unit struct {
	Entity
	Damage   float64
	Speed    float64
	AttackRange float64
	TargetX  float64
	TargetY  float64
	Moving   bool
	Attacking bool
	UnitType string  // "soldier", "tank", "ranger"
}

// Enemy representa um inimigo (wave PvE)
type Enemy struct {
	Entity
	Speed    float64
	Damage   float64
	Bounty   float64 // ouro ao morrer
	Wave     int
	PathIndex int
}

// Hero representa um herói
type Hero struct {
	Entity
	Damage          float64
	Speed           float64
	Level           int
	Experience      float64
	Ability         HeroAbility
	CanUseAbility   bool
	HeroType        string // "tank", "support", "damage", "swarm"
}

// SwarmUnit representa uma unidade do Enxame Voraz
type SwarmUnit struct {
	Entity
	Damage    float64
	Speed     float64
	AttackRange float64
	TargetX   float64
	TargetY   float64
	Moving    bool
	Attacking bool
	UnitType  string // "drone", "stalker", "behemoth"
	EssenceCost float64
}

// Projectile representa um projétil
type Projectile struct {
	Entity
	Damage    float64
	Speed     float64
	TargetID  uint64
	SourceID  uint64
}
