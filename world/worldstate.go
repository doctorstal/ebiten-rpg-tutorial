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
	AttackItems      []entities.AttackItem
}

func (r *RoomState) positionPlayerAtTheDoor(doorName string) {
	if door, ok := r.TiledMap.Doors()[doorName]; ok {
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

	room.Enemies = make([]*entities.Enemy, len(tiledMap.Enemies()))
	for i, e := range tiledMap.Enemies() {
		room.Enemies[i] = entities.NewEnemy(float64(e.Rect.Min.X), float64(e.Rect.Min.Y), e.FollorsPlayer, loader)
	}

	room.Potions = make([]*entities.Potion, 0)
	for _, item := range tiledMap.Items() {
		if item.Kind == "LifePotion" {
			room.Potions = append(room.Potions, entities.NewPotion(float64(item.Rect.Min.X), float64(item.Rect.Min.Y), loader))
		}
	}
	room.TiledMap = tiledMap
	colliders := make([]*image.Rectangle, 0)

	for _, objectRect := range room.TiledMap.ObjectRects() {
		colliders = append(colliders, objectRect)
	}
	room.Colliders = colliders
	room.AttackItems = make([]entities.AttackItem, 0)

	return room

}

type WorldState struct {
	CurrentRoom string
	rooms       map[string]*RoomState
	loader      *resource.Loader
}

func (w *WorldState) LoadRoom(roomFile string, player *entities.Player) *RoomState {
	if room, loaded := w.rooms[roomFile]; loaded {
		room.positionPlayerAtTheDoor(w.CurrentRoom)
		w.CurrentRoom = roomFile
		return room
	}
	room := NewRoom(roomFile, player, w.loader)
	room.positionPlayerAtTheDoor(w.CurrentRoom)
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
