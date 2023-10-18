package images

import (
	_ "embed"
)

var (
	//go:embed BOX.png
	BOX []byte

	//go:embed BOX_ON_GOAL.png
	BOX_ON_GOAL []byte

	//go:embed FLOOR.png
	FLOOR []byte

	//go:embed GOAL.png
	GOAL []byte

	//go:embed PLAYER.png
	PLAYER []byte

	//go:embed WALL.png
	WALL []byte
)
