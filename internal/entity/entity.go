package entity

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
}

// Unit representa uma unidade controlável
type Unit struct {
	Entity
	Damage   float64
	Speed    float64
	TargetX  float64
	TargetY  float64
	Moving   bool
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
	Damage     float64
	Speed      float64
	Level      int
	Experience float64
	AbilityCooldown float64
	CanUseAbility   bool
}

// Projectile representa um projétil
type Projectile struct {
	Entity
	Damage    float64
	Speed     float64
	TargetID  uint64
	SourceID  uint64
}
