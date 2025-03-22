package entities

import (
	"rpg-tutorial/components"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	*Character
	Bombs   []*Bomb
	BombImg *ebiten.Image
	Health  uint
}

func (p *Player) NewBomb() *Bomb {
	return NewBomb(p.BombImg, p.X, p.Y, p.CombatComponent.AttackPower())
}

func NewPlayer(playerImg, bombImg *ebiten.Image, x, y float64) *Player {
	return &Player{
		Character: NewCharacter(playerImg, x, y, components.NewPlayerCombat(5, 1, 10)),
		BombImg:   bombImg,
		Bombs:     make([]*Bomb, 0),
		Health:    3,
	}
}
