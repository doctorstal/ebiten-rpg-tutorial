package scenes

import (
	"rpg-tutorial/resources"
	"rpg-tutorial/state"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	resource "github.com/quasilyte/ebitengine-resource"
)

type StartScene struct {
	gameState *state.GlobalGameState
	loader    *resource.Loader
	selected  bool
	isLoaded  bool
	ui        ebitenui.UI
}

// Draw implements Scene.
func (s *StartScene) Draw(screen *ebiten.Image) {
	s.ui.Draw(screen)
	ebitenutil.DebugPrint(screen, "Choose your hero.")
}

// FirstLoad implements Scene.
func (s *StartScene) FirstLoad() {
	// load images for button states: idle, hover, and pressed

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSlice(s.loader.LoadImage(resources.UiPanelBg).Data, [3]int{6, 3, 6}, [3]int{6, 3, 6})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	innerContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
			widget.WidgetOpts.MinSize(32, 100),
		),
	)

	samuraiContainer, samuraiBtn := newButtonWithImage(s.loader, resources.ImgSamuraiFace, func() {
		s.selectCharacter(state.HeroSamurai)
	})

	robotContainer, robotBtn := newButtonWithImage(s.loader, resources.ImgRobotFace, func() {
		s.selectCharacter(state.HeroRobot)
	})

	skeletonContainer, skeletonBtn := newButtonWithImage(s.loader, resources.ImgSkeletonFace, func() {
		s.selectCharacter(state.HeroSkeleton)
	})

	samuraiBtn.AddFocus(widget.FOCUS_NEXT, robotBtn)
	robotBtn.AddFocus(widget.FOCUS_NEXT, skeletonBtn)
	skeletonBtn.AddFocus(widget.FOCUS_NEXT, samuraiBtn)
	// skeletonBtn.Focus(true)

	innerContainer.AddChild(samuraiContainer)
	innerContainer.AddChild(robotContainer)
	innerContainer.AddChild(skeletonContainer)

	rootContainer.AddChild(innerContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	s.ui = ui
	s.isLoaded = true
}
func (s *StartScene) selectCharacter(c state.Hero) {
	s.gameState.SelectedHero = c
	s.selected = true
}

func newButtonWithImage(loader *resource.Loader, spriteId resource.ImageID, callback func()) (*widget.Container, *widget.Button) {
	buttonImage := loadButtonImage(loader)
	buttonIcon := loader.LoadImage(spriteId).Data

	buttonStackedLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
		// instruct the container's anchor layout to center the button both horizontally and vertically;
		// since our button is a 2-widget object, we add the anchor info to the wrapping container
		// instead of the button
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	)
	btnIconG := widget.NewGraphic(
		widget.GraphicOpts.Images(&widget.GraphicImage{
			Idle: buttonIcon,
		}),
	)
	menuSound2 := loader.LoadAudio(resources.SoundMenu2).Player
	menuSound1 := loader.LoadAudio(resources.SoundMenu1).Player

	// construct a pressable button
	button := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			menuSound1.Rewind()
			menuSound1.Play()
			callback()
		}),
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			menuSound2.Rewind()
			menuSound2.Play()
		}),
	)
	button.GetWidget().FocusEvent.AddHandler(func(args any) {
			menuSound2.Rewind()
			menuSound2.Play()
	})
	buttonStackedLayout.AddChild(button)
	// Put an image on top of the button, it will be centered.
	// If your image doesn't fit the button and there is no Y stretching support,
	// you may see a transparent rectangle inside the button.
	// To fix that, either use a separate button image (that can fit the image)
	// or add an appropriate stretching.
	buttonStackedLayout.AddChild(
		btnIconG,
	)
	return buttonStackedLayout, button
}

// IsLoaded implements Scene.
func (s *StartScene) IsLoaded() bool {
	return s.isLoaded
}

// OnEnter implements Scene.
func (s *StartScene) OnEnter() {
	s.selected = false
}

// OnExit implements Scene.
func (s *StartScene) OnExit() {

}

// Update implements Scene.
func (s *StartScene) Update() SceneId {
	s.ui.Update()
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ExitSceneId
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		s.ui.ChangeFocus(widget.FOCUS_PREVIOUS)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		s.ui.ChangeFocus(widget.FOCUS_NEXT)
	}
	if s.selected {
		s.selected = false
		s.ui.ClearFocus()
		return GameSceneId
	}
	return StartSceneId
}

func NewStartScene(gameState *state.GlobalGameState, loader *resource.Loader) Scene {
	return &StartScene{gameState: gameState, loader: loader}
}

func loadButtonImage(loader *resource.Loader) *widget.ButtonImage {
	idle := eimage.NewNineSlice(loader.LoadImage(resources.UiBtnNormal).Data, [3]int{3, 10, 3}, [3]int{3, 3, 3})
	hover := eimage.NewNineSlice(loader.LoadImage(resources.UiBtnHover).Data, [3]int{3, 10, 3}, [3]int{3, 3, 3})
	pressed := eimage.NewNineSlice(loader.LoadImage(resources.UiBtnPressed).Data, [3]int{3, 10, 3}, [3]int{3, 3, 3})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}
