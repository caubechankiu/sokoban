package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"time"

	"github.com/caubechankiu/sokoban/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/samber/lo"
)

const SCREEN_WIDTH int = 640
const SCREEN_HEIGHT int = 480

type LevelElement string

const WALL LevelElement = "#"
const FLOOR LevelElement = " "
const BOX LevelElement = "$"
const BOX_ON_GOAL LevelElement = "*"
const PLAYER LevelElement = "@"
const PLAYER_ON_GOAL LevelElement = "+"
const GOAL LevelElement = "."

var (
	floorImage     *ebiten.Image
	wallImage      *ebiten.Image
	playerImage    *ebiten.Image
	boxImage       *ebiten.Image
	boxOnGoalImage *ebiten.Image
	goalImage      *ebiten.Image
)

func init() {
	box, _, err := image.Decode(bytes.NewReader(images.BOX))
	if err != nil {
		log.Fatal(err)
	}
	boxImage = ebiten.NewImageFromImage(box)

	floor, _, err := image.Decode(bytes.NewReader(images.FLOOR))
	if err != nil {
		log.Fatal(err)
	}
	floorImage = ebiten.NewImageFromImage(floor)

	goal, _, err := image.Decode(bytes.NewReader(images.GOAL))
	if err != nil {
		log.Fatal(err)
	}
	goalImage = ebiten.NewImageFromImage(goal)

	player, _, err := image.Decode(bytes.NewReader(images.PLAYER))
	if err != nil {
		log.Fatal(err)
	}
	playerImage = ebiten.NewImageFromImage(player)

	boxOnGoal, _, err := image.Decode(bytes.NewReader(images.BOX_ON_GOAL))
	if err != nil {
		log.Fatal(err)
	}
	boxOnGoalImage = ebiten.NewImageFromImage(boxOnGoal)

	wall, _, err := image.Decode(bytes.NewReader(images.WALL))
	if err != nil {
		log.Fatal(err)
	}
	wallImage = ebiten.NewImageFromImage(wall)
}

type Game struct {
	Map              [][]LevelElement
	MapSizeX         int
	MapSizeY         int
	PlayerX          int
	PlayerY          int
	IsPlayerOnGoal   bool
	LastTimeKeyPress int64
}

func (g *Game) Update() error {
	keys := inpututil.AppendPressedKeys([]ebiten.Key{})
	if len(keys) > 0 && g.LastTimeKeyPress < time.Now().UnixMilli()-200 {
		g.LastTimeKeyPress = time.Now().UnixMilli()
		switch keys[0] {
		case ebiten.KeyA:
			log.Println(g.PlayerX, g.PlayerY)
		case ebiten.KeyArrowUp:
			if g.Map[g.PlayerY-1][g.PlayerX] == FLOOR {
				g.Map[g.PlayerY-1][g.PlayerX] = PLAYER
				g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
				g.PlayerY = g.PlayerY - 1
				g.IsPlayerOnGoal = false
				break
			} else if g.Map[g.PlayerY-1][g.PlayerX] == BOX {
				if g.Map[g.PlayerY-2][g.PlayerX] == FLOOR {
					g.Map[g.PlayerY-2][g.PlayerX] = BOX
					g.Map[g.PlayerY-1][g.PlayerX] = PLAYER
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerY = g.PlayerY - 1
					g.IsPlayerOnGoal = false
				} else if g.Map[g.PlayerY-2][g.PlayerX] == GOAL {
					g.Map[g.PlayerY-2][g.PlayerX] = BOX_ON_GOAL
					g.Map[g.PlayerY-1][g.PlayerX] = PLAYER
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerY = g.PlayerY - 1
					g.IsPlayerOnGoal = false
				}
			} else if g.Map[g.PlayerY-1][g.PlayerX] == BOX_ON_GOAL {
				if g.Map[g.PlayerY-2][g.PlayerX] == FLOOR {
					g.Map[g.PlayerY-2][g.PlayerX] = BOX
					g.Map[g.PlayerY-1][g.PlayerX] = PLAYER_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerY = g.PlayerY - 1
					g.IsPlayerOnGoal = true
				} else if g.Map[g.PlayerY-2][g.PlayerX] == GOAL {
					g.Map[g.PlayerY-2][g.PlayerX] = BOX_ON_GOAL
					g.Map[g.PlayerY-1][g.PlayerX] = PLAYER_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerY = g.PlayerY - 1
					g.IsPlayerOnGoal = true
				}
			} else if g.Map[g.PlayerY-1][g.PlayerX] == GOAL {
				g.Map[g.PlayerY-1][g.PlayerX] = PLAYER_ON_GOAL
				g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
				g.PlayerY = g.PlayerY - 1
				g.IsPlayerOnGoal = true
				break
			}
		case ebiten.KeyArrowDown:
			if g.Map[g.PlayerY+1][g.PlayerX] == FLOOR {
				g.Map[g.PlayerY+1][g.PlayerX] = PLAYER
				g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
				g.PlayerY = g.PlayerY + 1
				g.IsPlayerOnGoal = false
				break
			} else if g.Map[g.PlayerY+1][g.PlayerX] == BOX {
				if g.Map[g.PlayerY+2][g.PlayerX] == FLOOR {
					g.Map[g.PlayerY+2][g.PlayerX] = BOX
					g.Map[g.PlayerY+1][g.PlayerX] = PLAYER
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerY = g.PlayerY + 1
					g.IsPlayerOnGoal = false
				} else if g.Map[g.PlayerY+2][g.PlayerX] == GOAL {
					g.Map[g.PlayerY+2][g.PlayerX] = BOX_ON_GOAL
					g.Map[g.PlayerY+1][g.PlayerX] = PLAYER
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerY = g.PlayerY + 1
					g.IsPlayerOnGoal = false
				}
			} else if g.Map[g.PlayerY+1][g.PlayerX] == BOX_ON_GOAL {
				if g.Map[g.PlayerY+2][g.PlayerX] == FLOOR {
					g.Map[g.PlayerY+2][g.PlayerX] = BOX
					g.Map[g.PlayerY+1][g.PlayerX] = PLAYER_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerY = g.PlayerY + 1
					g.IsPlayerOnGoal = true
				} else if g.Map[g.PlayerY+2][g.PlayerX] == GOAL {
					g.Map[g.PlayerY+2][g.PlayerX] = BOX_ON_GOAL
					g.Map[g.PlayerY+1][g.PlayerX] = PLAYER_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerY = g.PlayerY + 1
					g.IsPlayerOnGoal = true
				}
			} else if g.Map[g.PlayerY+1][g.PlayerX] == GOAL {
				g.Map[g.PlayerY+1][g.PlayerX] = PLAYER_ON_GOAL
				g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
				g.PlayerY = g.PlayerY + 1
				g.IsPlayerOnGoal = true
				break
			}
		case ebiten.KeyArrowLeft:
			if g.Map[g.PlayerY][g.PlayerX-1] == FLOOR {
				g.Map[g.PlayerY][g.PlayerX-1] = PLAYER
				g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
				g.PlayerX = g.PlayerX - 1
				g.IsPlayerOnGoal = false
				break
			} else if g.Map[g.PlayerY][g.PlayerX-1] == BOX {
				if g.Map[g.PlayerY][g.PlayerX-2] == FLOOR {
					g.Map[g.PlayerY][g.PlayerX-2] = BOX
					g.Map[g.PlayerY][g.PlayerX-1] = PLAYER
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerX = g.PlayerX - 1
					g.IsPlayerOnGoal = false
				} else if g.Map[g.PlayerY][g.PlayerX-2] == GOAL {
					g.Map[g.PlayerY][g.PlayerX-2] = BOX_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX-1] = PLAYER
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerX = g.PlayerX - 1
					g.IsPlayerOnGoal = false
				}
			} else if g.Map[g.PlayerY][g.PlayerX-1] == BOX_ON_GOAL {
				if g.Map[g.PlayerY][g.PlayerX-2] == FLOOR {
					g.Map[g.PlayerY][g.PlayerX-2] = BOX
					g.Map[g.PlayerY][g.PlayerX-1] = PLAYER_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerX = g.PlayerX - 1
					g.IsPlayerOnGoal = true
				} else if g.Map[g.PlayerY][g.PlayerX-2] == GOAL {
					g.Map[g.PlayerY][g.PlayerX-2] = BOX_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX-1] = PLAYER_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerX = g.PlayerX - 1
					g.IsPlayerOnGoal = true
				}
			} else if g.Map[g.PlayerY][g.PlayerX-1] == GOAL {
				g.Map[g.PlayerY][g.PlayerX-1] = PLAYER_ON_GOAL
				g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
				g.PlayerX = g.PlayerX - 1
				g.IsPlayerOnGoal = true
				break
			}
		case ebiten.KeyArrowRight:
			if g.Map[g.PlayerY][g.PlayerX+1] == FLOOR {
				g.Map[g.PlayerY][g.PlayerX+1] = PLAYER
				g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
				g.PlayerX = g.PlayerX + 1
				g.IsPlayerOnGoal = false
				break
			} else if g.Map[g.PlayerY][g.PlayerX+1] == BOX {
				if g.Map[g.PlayerY][g.PlayerX+2] == FLOOR {
					g.Map[g.PlayerY][g.PlayerX+2] = BOX
					g.Map[g.PlayerY][g.PlayerX+1] = PLAYER
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerX = g.PlayerX + 1
					g.IsPlayerOnGoal = false
				} else if g.Map[g.PlayerY][g.PlayerX+2] == GOAL {
					g.Map[g.PlayerY][g.PlayerX+2] = BOX_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX+1] = PLAYER
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerX = g.PlayerX + 1
					g.IsPlayerOnGoal = false
				}
			} else if g.Map[g.PlayerY][g.PlayerX+1] == BOX_ON_GOAL {
				if g.Map[g.PlayerY][g.PlayerX+2] == FLOOR {
					g.Map[g.PlayerY][g.PlayerX+2] = BOX
					g.Map[g.PlayerY][g.PlayerX+1] = PLAYER_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerX = g.PlayerX + 1
					g.IsPlayerOnGoal = true
				} else if g.Map[g.PlayerY][g.PlayerX+2] == GOAL {
					g.Map[g.PlayerY][g.PlayerX+2] = BOX_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX+1] = PLAYER_ON_GOAL
					g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
					g.PlayerX = g.PlayerX + 1
					g.IsPlayerOnGoal = true
				}
			} else if g.Map[g.PlayerY][g.PlayerX+1] == GOAL {
				g.Map[g.PlayerY][g.PlayerX+1] = PLAYER_ON_GOAL
				g.Map[g.PlayerY][g.PlayerX] = lo.If[LevelElement](g.IsPlayerOnGoal, GOAL).Else(FLOOR)
				g.PlayerX = g.PlayerX + 1
				g.IsPlayerOnGoal = true
				break
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y, row := range g.Map {
		for x, element := range row {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.25, 0.25)
			op.GeoM.Translate(float64(SCREEN_WIDTH/2-32*g.MapSizeX/2+32*x), float64(SCREEN_HEIGHT/2-32*g.MapSizeY/2+32*y))
			switch LevelElement(element) {
			case FLOOR:
				screen.DrawImage(floorImage, op)
			case BOX:
				screen.DrawImage(boxImage, op)
			case BOX_ON_GOAL:
				screen.DrawImage(boxOnGoalImage, op)
			case PLAYER:
				screen.DrawImage(floorImage, op)
				screen.DrawImage(playerImage, op)
			case PLAYER_ON_GOAL:
				screen.DrawImage(floorImage, op)
				screen.DrawImage(playerImage, op)
			case WALL:
				screen.DrawImage(wallImage, op)
			case GOAL:
				screen.DrawImage(floorImage, op)
				screen.DrawImage(goalImage, op)
			}
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nPlayerX: %d, PlayerY: %d, IsPlayerOnGoal: %t, MapSizeX: %d, MapSizeY: %d", g.PlayerX, g.PlayerY, g.IsPlayerOnGoal, g.MapSizeX, g.MapSizeY))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = SCREEN_WIDTH
	screenHeight = SCREEN_HEIGHT
	return
}

func main() {
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_WIDTH)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Sokoban")

	game := &Game{Map: [][]LevelElement{}}

	if err := game.LoadLevel("levels/picokosmos/13.txt"); err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) LoadLevel(filePath string) error {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	row := []LevelElement{}
	y := 0
	x := 0
	for _, char := range string(buf) {
		if string(char) == "\n" {
			g.Map = append(g.Map, row)
			row = []LevelElement{}
			y++
			x = 0
			g.MapSizeY = y
		} else {
			row = append(row, LevelElement(char))
			if LevelElement(char) == PLAYER || LevelElement(char) == PLAYER_ON_GOAL {
				g.PlayerX = x
				g.PlayerY = y
			}
			if LevelElement(char) == PLAYER_ON_GOAL {
				g.IsPlayerOnGoal = true
			}
			x++
			g.MapSizeX = int(math.Max(float64(g.MapSizeX), float64(x)))
		}
	}
	if len(row) > 0 {
		g.Map = append(g.Map, row)
		g.MapSizeY = y + 1
	}

	return nil
}
