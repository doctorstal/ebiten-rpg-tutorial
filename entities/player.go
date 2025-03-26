package entities

import (
	"fmt"
	"log"
	"rpg-tutorial/components"
	"rpg-tutorial/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	*Character
	AttackItems []AttackItem
	BombImg     *ebiten.Image
	hero        state.Hero
}

func (p *Player) NewBomb() AttackItem {
	switch p.hero {
	case state.HeroRobot:
		return NewBomb(p.BombImg, p.X, p.Y, p.CombatComponent.AttackPower())
	case state.HeroSamurai:
		return NewEnergyBall(p.BombImg, p.X, p.Y, p.CombatComponent.AttackPower(), p.Direction, 8.0)
	case state.HeroSkeleton:
		return NewRock(p.BombImg, p.X, p.Y, p.CombatComponent.AttackPower(), p.Direction, 5.0)
	default:
		panic(fmt.Sprintf("unexpected state.Hero: %#v", p.hero))
	}
}

func NewPlayer(hero state.Hero, x, y float64) *Player {
	var playerImgPath = "assets/images/samurai.png"
	var bombImgPath = "assets/images/weapons/bomb.png"

	switch hero {
	case state.HeroRobot:
		playerImgPath = "assets/images/robot.png"
	case state.HeroSamurai:
		playerImgPath = "assets/images/samurai.png"
		bombImgPath = "assets/images/weapons/energy_ball.png"
	case state.HeroSkeleton:
		playerImgPath = "assets/images/skeleton.png"
		bombImgPath = "assets/images/weapons/rock.png"
	default:
		panic(fmt.Sprintf("unexpected state.Hero: %#v", hero))
	}

	playerImg, _, err := ebitenutil.NewImageFromFile(playerImgPath)
	if err != nil {
		log.Fatal(err)
	}

	bombImg, _, err := ebitenutil.NewImageFromFile(bombImgPath)
	if err != nil {
		log.Fatal(err)
	}
	return &Player{
		Character:   NewCharacter(playerImg, x, y, components.NewPlayerCombat(5, 1, 30, 10)),
		BombImg:     bombImg,
		AttackItems: make([]AttackItem, 0),
		hero:        hero,
	}
}
