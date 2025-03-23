package components

type Combat interface {
	Health() int
	AttackPower() uint
	Attacking() bool
	Attack() bool
	Update()
	Damage(amount uint)
	Heal(amount uint)
}

type BasicCombat struct {
	health      int
	attackPower uint
	attacking   bool
}

// Attack implements Combat.
func (b *BasicCombat) Attack() bool {
	b.attacking = true
	return true
}

func (b *BasicCombat) Update() {
	b.attacking = false
}

// Attacking implements Combat.
func (b *BasicCombat) Attacking() bool {
	return b.attacking
}

func (b *BasicCombat) AttackPower() uint {
	return b.attackPower
}

func (b *BasicCombat) Damage(amount uint) {
	b.health -= int(amount)
}

func (b *BasicCombat) Heal(amount uint) {
	b.health += int(amount)
}

func (b *BasicCombat) Health() int {
	return b.health
}
func NewBasicCombat(health int, attackPower uint) Combat {
	return &BasicCombat{
		health:      health,
		attackPower: attackPower,
		attacking:   false,
	}
}

type PlayerCombat struct {
	*BasicCombat
	attackCooldown  int
	attackImmobile  int
	timeSinceAttack int
}

func NewPlayerCombat(health int, attackPower uint, attackCooldown int, attackImmobile int) Combat {
	return &PlayerCombat{
		BasicCombat:     NewBasicCombat(health, attackPower).(*BasicCombat),
		attackCooldown:  attackCooldown,
		attackImmobile:  attackImmobile,
		timeSinceAttack: 1000,
	}
}

func (p *PlayerCombat) Attack() bool {
	if p.timeSinceAttack >= p.attackCooldown {
		p.attacking = true
		p.timeSinceAttack = 0
		return true
	}
	return false
}

func (p *PlayerCombat) Update() {
	if p.timeSinceAttack < 1000 {
		p.timeSinceAttack++
	}
	if p.timeSinceAttack >= p.attackImmobile {
		p.attacking = false
	}
}

type EnemyCombat struct {
	*BasicCombat
	attackCooldown  int
	timeSinceAttack int
}

func NewEnemyCombat(health int, attackPower uint, attackCooldown int) Combat {
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
	if e.timeSinceAttack < 1000 {
		e.timeSinceAttack++
	}
	if e.timeSinceAttack >= e.attackCooldown {
		e.attacking = false
	}
}
