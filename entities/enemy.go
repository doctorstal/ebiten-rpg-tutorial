package entities

import (
	"image"

	"github.com/doctorstal/ebiten-rpg-tutorial/components"
	"github.com/doctorstal/ebiten-rpg-tutorial/resources"
	"github.com/hajimehoshi/ebiten/v2"

	resource "github.com/quasilyte/ebitengine-resource"
)

type LifeBar struct {
	img, imgUnder, renderedImg *ebiten.Image
}

func (l *LifeBar) GetRenderer(health, totalHealth int, rect *image.Rectangle) Renderer {
	opts := &ebiten.DrawImageOptions{}
	l.renderedImg.DrawImage(l.imgUnder, opts)
	l.renderedImg.DrawImage(
		l.img.SubImage(image.Rect(
			0,
			0,
			l.img.Bounds().Dx()*health/totalHealth,
			l.img.Bounds().Dy(),
		)).(*ebiten.Image),
		opts,
	)

	movedRect := rect.Add(image.Point{0, -5})
	return &BasicRenderer{
		img:      l.renderedImg,
		drawOpts: &ebiten.DrawImageOptions{},
		rect:     &movedRect,
		z:        1, // On top of everything
	}
}

func NewLifeBar(loader *resource.Loader) *LifeBar {
	lifebarImg := loader.LoadImage(resources.ImgLifeBar).Data
	return &LifeBar{
		img:         lifebarImg,
		imgUnder:    loader.LoadImage(resources.ImgLifeBarUnder).Data,
		renderedImg: ebiten.NewImage(lifebarImg.Bounds().Dx(), lifebarImg.Bounds().Dy()),
	}
}

type Enemy struct {
	*Character
	FollowsPlayer  bool
	WonderingSpeed float64
	lifebar        *LifeBar
}

func (e *Enemy) GetRenderers() []Renderer {
	lifebar := e.lifebar.GetRenderer(
		e.CombatComponent.Health(),
		e.CombatComponent.MaxHealth(),
		e.Rect(),
	)
	return append(e.Sprite.GetRenderers(), lifebar)
}

func (e *Enemy) IsDead() bool {
	return e.state == Dead
}

func NewEnemy(x, y float64, fp bool, loader *resource.Loader) *Enemy {
	return &Enemy{
		Character: NewCharacter(
			loader.LoadImage(resources.ImgSkeleton).Data,
			x,
			y,
			components.NewEnemyCombat(3, 1, 30),
		),
		lifebar:       NewLifeBar(loader),
		FollowsPlayer: fp,
	}
}
