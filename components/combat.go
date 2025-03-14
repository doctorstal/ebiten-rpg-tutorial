package components


type Combat interface {
	Health() int
	AttackPower() int
	Attacking() bool
	Attack() bool
	Update()
	Damage(amount int)
}

type BasicCombat struct {
	health      int
	attackPower int
	attacking   bool
}

// Attack implements Combat.
func (b *BasicCombat) Attack() bool {
	b.attacking = true
	return true
}

func (b *BasicCombat) Update() {

}

// Attacking implements Combat.
func (b *BasicCombat) Attacking() bool {
	return b.attacking
}

func (b *BasicCombat) AttackPower() int {
	return b.attackPower
}

func (b *BasicCombat) Damage(amount int) {
	b.health -= amount
}

func (b *BasicCombat) Health() int {
	return b.health
}
func NewBasicCombat(health, attackPower int) Combat {
	return &BasicCombat{
		health:      health,
		attackPower: attackPower,
		attacking:   false,
	}
}

type EnemyCombat struct {
	*BasicCombat
	attackCooldown  int
	timeSinceAttack int
}

func NewEnemyCombat(health, attackPower, attackCooldown int) Combat {
	return &EnemyCombat{
		BasicCombat:     NewBasicCombat(health, attackPower).(*BasicCombat),
		attackCooldown:  attackCooldown,
		timeSinceAttack: 0,
	}
}

func (e *EnemyCombat) Attack() bool {
	if e.timeSinceAttack >= e.attackCooldown {
		e.attacking = true
		e.timeSinceAttack = 0
		return true
	}
	return false
}

func (e *EnemyCombat) Update() {
	e.timeSinceAttack++
}
