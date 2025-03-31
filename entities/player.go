package entities

import (
	"fmt"
	"github.com/doctorstal/ebiten-rpg-tutorial/components"
	"github.com/doctorstal/ebiten-rpg-tutorial/resources"
	"github.com/doctorstal/ebiten-rpg-tutorial/state"

	resource "github.com/quasilyte/ebitengine-resource"
)

type Player struct {
	*Character
	hero        state.Hero
	loader      *resource.Loader
}

func (p *Player) NewAttackItem() AttackItem {
	switch p.hero {
	case state.HeroRobot:
		return NewBomb(p.loader, p.X, p.Y+1, p.CombatComponent.AttackPower())
	case state.HeroSamurai:
		return NewEnergyBall(p.loader, p.X, p.Y, p.CombatComponent.AttackPower(), p.Direction, 8.0)
	case state.HeroBoy:
		return NewRock(p.loader, p.X, p.Y, p.CombatComponent.AttackPower(), p.Direction, 5.0)
	default:
		panic(fmt.Sprintf("unexpected state.Hero: %#v", p.hero))
	}
}

func NewPlayer(loader *resource.Loader, hero state.Hero, x, y float64) *Player {
	playerImgId := resources.ImgSamurai

	switch hero {
	case state.HeroRobot:
		playerImgId = resources.ImgRobot
	case state.HeroSamurai:
		playerImgId = resources.ImgSamurai
	case state.HeroBoy:
		playerImgId = resources.ImgBoy
	default:
		panic(fmt.Sprintf("unexpected state.Hero: %#v", hero))
	}

	playerImg := loader.LoadImage(playerImgId).Data
	return &Player{
		Character:   NewCharacter(playerImg, x, y, components.NewPlayerCombat(5, 1, 30, 10)),
		hero:        hero,
		loader:      loader,
	}
}
