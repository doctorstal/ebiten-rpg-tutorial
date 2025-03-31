package world

import (
	"image"
	"log"
	"path"
	"github.com/doctorstal/ebiten-rpg-tutorial/entities"
	"github.com/doctorstal/ebiten-rpg-tutorial/tiled"

	resource "github.com/quasilyte/ebitengine-resource"
)

type RoomState struct {
	Player           *entities.Player
	playerX, playerY float64 // For switching between levels
	Enemies          []*entities.Enemy
	StaticAnimators  []entities.Animator
	Potions          []*entities.Potion
	TiledMap         *tiled.TiledMap
	Colliders        []*image.Rectangle
}

func (r *RoomState) positionPlayer(from string) {
	if door, ok := r.TiledMap.Doors()[from]; ok {
		r.Player.X = float64(door.Rect.Min.X)
		if door.Direction == "up" {
			r.Player.Y = float64(door.Rect.Min.Y) - r.Player.Height
		} else {
			r.Player.Y = float64(door.Rect.Max.Y)
		}
	} else {
		r.Player.X = 0
		r.Player.Y = 0
	}

}

func NewRoom(mapFile string, player *entities.Player, loader *resource.Loader) *RoomState {
	tiledMap, err := tiled.NewTiledMap(path.Join("assets", "maps", mapFile))
	if err != nil {
		log.Fatal(err)
	}

	room := &RoomState{}

	room.Player = player

	room.Enemies = []*entities.Enemy{
		entities.NewEnemy(50.0, 50.0, false, loader),
		entities.NewEnemy(75.0, 75.0, false, loader),
		entities.NewEnemy(150.0, 75.0, false, loader),
		// entities.NewEnemy(150.0, 75.0, true, loader),
		// entities.NewEnemy(150.0, 75.0, true, loader),
		// entities.NewEnemy(150.0, 75.0, false, loader),
		// entities.NewEnemy(150.0, 75.0, true, loader),
		// entities.NewEnemy(150.0, 75.0, true, loader),
		// entities.NewEnemy(150.0, 75.0, true, loader),
		// entities.NewEnemy(150.0, 75.0, true, loader),
		// entities.NewEnemy(150.0, 75.0, true, loader),
		// entities.NewEnemy(150.0, 75.0, true, loader),
	}
	room.Potions = []*entities.Potion{
		entities.NewPotion(210.0, 100.0, loader),
	}

	room.TiledMap = tiledMap
	colliders := make([]*image.Rectangle, 0)

	for _, objectRect := range room.TiledMap.ObjectRects() {
		colliders = append(colliders, objectRect)
	}
	room.Colliders = colliders

	return room

}

type WorldState struct {
	CurrentRoom string
	rooms       map[string]*RoomState
	loader      *resource.Loader
}

func (w *WorldState) LoadRoom(roomFile string, player *entities.Player) *RoomState {
	if room, loaded := w.rooms[roomFile]; loaded {
		room.positionPlayer(w.CurrentRoom)
		w.CurrentRoom = roomFile
		return room
	}
	room := NewRoom(roomFile, player, w.loader)
	room.positionPlayer(w.CurrentRoom)
	w.CurrentRoom = roomFile
	w.rooms[roomFile] = room
	return room
}

func NewWorldState(loader *resource.Loader, currentRoom string) *WorldState {
	return &WorldState{
		CurrentRoom: currentRoom,
		rooms:       make(map[string]*RoomState),
		loader:      loader,
	}
}
