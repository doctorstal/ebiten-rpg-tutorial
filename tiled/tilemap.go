package tiled

import (
	"encoding/json"
	"os"
	"path"
)

type TilemapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `jaon:"name"`
}
type TilemapJSON struct {
	Layers   []*TilemapLayerJSON `json:"layers"`
	Tilesets []map[string]any    `json:"tilesets"`
}

func (t *TilemapJSON) GenTilesets() ([]Tileset, error) {
	tilesets := make([]Tileset, 0)
	for _, tilesetsData := range t.Tilesets {
		tilesetPath := path.Join("assets", "maps", tilesetsData["source"].(string))
		tilesetGid := int(tilesetsData["firstgid"].(float64))
		tileset, err := NewTileset(tilesetPath, tilesetGid)
		if err != nil {
			return nil, err
		}
		tilesets = append(tilesets, tileset)

	}

	return tilesets, nil
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
