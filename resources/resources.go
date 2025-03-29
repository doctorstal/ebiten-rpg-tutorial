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
	ImgBombExplosion
	ImgEnergyBall
	ImgRock
	ImgRockExplosion
	ImgShadow
	ImgPotion
	ImgHeart

	UiBtnNormal
	UiBtnHover
	UiBtnPressed
	UiPanelBg
)

const (
	SoundMenu1 resource.AudioID = iota
	SoundMenu2
	SoundExplosion
	SoundFireball
	SoundRockSmash

	SoundLost
	SoundWon
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
		ImgSkeleton:      {Path: "assets/images/skeleton.png"},
		ImgSkeletonFace:  {Path: "assets/images/skeleton_faceset.png"},
		ImgSamurai:       {Path: "assets/images/samurai.png"},
		ImgSamuraiFace:   {Path: "assets/images/samurai_faceset.png"},
		ImgRobot:         {Path: "assets/images/robot.png"},
		ImgRobotFace:     {Path: "assets/images/robot_faceset.png"},
		ImgBomb:          {Path: "assets/images/weapons/bomb.png"},
		ImgBombExplosion: {Path: "assets/images/weapons/explosion32x32.png"},
		ImgEnergyBall:    {Path: "assets/images/weapons/energy_ball.png"},
		ImgRock:          {Path: "assets/images/weapons/rock.png"},
		ImgRockExplosion: {Path: "assets/images/weapons/rock_element32x32.png"},
		ImgShadow:        {Path: "assets/images/shadow.png"},
		ImgPotion:        {Path: "assets/images/potion.png"},
		ImgHeart:        {Path: "assets/images/ui/heart.png"},

		UiBtnNormal:  {Path: "assets/images/ui/button_normal.png"},
		UiBtnHover:   {Path: "assets/images/ui/button_hover.png"},
		UiBtnPressed: {Path: "assets/images/ui/button_pressed.png"},
		UiPanelBg:    {Path: "assets/images/ui/nine_path_panel.png"},
	})

	l.AudioRegistry.Assign(map[resource.AudioID]resource.AudioInfo{
		SoundMenu1:     {Path: "assets/sounds/ui/Menu1.wav"},
		SoundMenu2:     {Path: "assets/sounds/ui/Menu2.wav"},
		SoundExplosion: {Path: "assets/sounds/game/Explosion.wav"},
		SoundRockSmash: {Path: "assets/sounds/game/rock_smash.wav"},
		SoundFireball:  {Path: "assets/sounds/game/Fireball.wav"},
		SoundLost:      {Path: "assets/sounds/ui/GameOver2.wav"},
		SoundWon:       {Path: "assets/sounds/ui/Success3.wav"},
	})

	return l
}
