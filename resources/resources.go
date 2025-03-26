package resources

import (
	"embed"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
)

const (
	ImgSkeleton resource.ImageID = iota
	ImgSkeletonFace
	ImgSamurai
	ImgSamuraiFace
	ImgRobot
	ImgRobotFace
	ImgBomb
	ImgEnergyBall
	ImgRock
	ImgShadow
	ImgPotion

	UiBtnNormal
	UiBtnHover
	UiBtnPressed
	UiPanelBg
)

func NewResourceLoader(fs embed.FS, audioContext *audio.Context) *resource.Loader {
	l := resource.NewLoader(audioContext)
	l.OpenAssetFunc = func(path string) io.ReadCloser {
		file, err := fs.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		return file
	}

	l.ImageRegistry.Assign(map[resource.ImageID]resource.ImageInfo{
		ImgSkeleton:     {Path: "assets/images/skeleton.png"},
		ImgSkeletonFace: {Path: "assets/images/skeleton_faceset.png"},
		ImgSamurai:      {Path: "assets/images/samurai.png"},
		ImgSamuraiFace:  {Path: "assets/images/samurai_faceset.png"},
		ImgRobot:        {Path: "assets/images/robot.png"},
		ImgRobotFace:    {Path: "assets/images/robot_faceset.png"},
		ImgBomb:         {Path: "assets/images/weapons/bomb.png"},
		ImgEnergyBall:   {Path: "assets/images/weapons/energy_ball.png"},
		ImgRock:         {Path: "assets/images/weapons/rock.png"},
		ImgShadow:       {Path: "assets/images/shadow.png"},
		ImgPotion:       {Path: "assets/images/potion.png"},

		UiBtnNormal:     {Path: "assets/images/ui/button_normal.png"},
		UiBtnHover:      {Path: "assets/images/ui/button_hover.png"},
		UiBtnPressed:    {Path: "assets/images/ui/button_pressed.png"},
		UiPanelBg:       {Path: "assets/images/ui/nine_path_panel.png"},
	})

	return l
}
