package scenes

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"rpg-tutorial/state"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type StartScene struct {
	gameState *state.GlobalGameState
	selected  bool
	isLoaded  bool
	ui        ebitenui.UI
}

// Draw implements Scene.
func (s *StartScene) Draw(screen *ebiten.Image) {
	s.ui.Draw(screen)
	ebitenutil.DebugPrint(screen, "Press <Enter> to start.")
}

// FirstLoad implements Scene.
func (s *StartScene) FirstLoad() {
	// load images for button states: idle, hover, and pressed

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSliceFromFile15x15("assets/images/ui/nine_path_panel.png")),
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

	samuraiContainer, samuraiBtn := newButtonWithImage("assets/images/samurai_faceset.png", func() {
		s.selectCharacter(state.HeroSamurai)
	})

	robotContainer, robotBtn := newButtonWithImage("assets/images/robot_faceset.png", func() {
		s.selectCharacter(state.HeroRobot)
	})

	skeletonContainer, skeletonBtn := newButtonWithImage("assets/images/skeleton_faceset.png", func() {
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

func newButtonWithImage(spritePath string, callback func()) (*widget.Container, *widget.Button) {
	buttonImage := loadButtonImage()
	buttonIcon := loadButtonIcon(spritePath)
	buttonDisabledIcon := loadDisabledButtonIcon()

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
			Idle:     buttonIcon,
			Disabled: buttonDisabledIcon,
		},
		),
	)
	// construct a pressable button
	button := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			callback()
		}),
	)
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		s.ui.ChangeFocus(widget.FOCUS_PREVIOUS)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {

		fmt.Println(s.ui.GetFocusedWidget())
		s.ui.ChangeFocus(widget.FOCUS_NEXT)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if btn, ok := s.ui.GetFocusedWidget().(*widget.Button); ok {
			fmt.Println(btn)
			// btn.Click()
		}
	}
	if s.selected {
		s.selected = false
		s.ui.ClearFocus()
		return GameSceneId
	}
	return StartSceneId
}

func NewStartScene(gameState *state.GlobalGameState) Scene {
	return &StartScene{gameState: gameState}
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
func loadButtonIcon(spritePath string) *ebiten.Image {
	icon, _, err := ebitenutil.NewImageFromFile(spritePath)
	if err != nil {
		log.Fatal(err)
	}
	return icon
}

func loadDisabledButtonIcon() *ebiten.Image {
	// we'll use a circle as an icon image
	// in reality it could be an arbitrary *ebiten.Image
	icon := ebiten.NewImage(32, 32)
	ebitenutil.DrawCircle(icon, 16, 16, 16, color.RGBA{R: 250, G: 0x56, B: 0xbd, A: 255})
	return icon
}

func loadButtonImage() *widget.ButtonImage {
	idle := nineSliceFromFile("assets/images/ui/button_normal.png")
	hover := nineSliceFromFile("assets/images/ui/button_hover.png")
	pressed := nineSliceFromFile("assets/images/ui/button_pressed.png")

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}
func nineSliceFromFile15x15(path string) *eimage.NineSlice {
	idleImg, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	idle := eimage.NewNineSlice(idleImg, [3]int{6, 3, 6}, [3]int{6, 3, 6})
	return idle
}
func nineSliceFromFile(path string) *eimage.NineSlice {
	idleImg, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	idle := eimage.NewNineSlice(idleImg, [3]int{3, 10, 3}, [3]int{3, 3, 3})
	return idle
}
