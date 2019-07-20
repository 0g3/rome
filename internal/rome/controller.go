package rome

import (
	"errors"
	"fmt"
)

type MoveDirection int

const (
	MoveDirectionUnknown MoveDirection = iota
	MoveDirectionLeft
	MoveDirectionRight
	MoveDirectionUp
	MoveDirectionDown
)

type GopherController struct {
	x, y   int
	bx, by int
}

func (gc *GopherController) Move(m Maze, md MoveDirection) error {
	var dx, dy int
	switch md {
	case MoveDirectionLeft:
		dx = -1
	case MoveDirectionRight:
		dx = 1
	case MoveDirectionUp:
		dy = 1
	case MoveDirectionDown:
		dy = -1
	default:
		return errors.New("不明なDirection")
	}

	newX := gc.x + dx
	newY := gc.y + dy
	if t, ok := m.Get(newX, newY); !ok || t != TileRoad && t != TileGoal {
		return fmt.Errorf("移動先が不正: (x,y)=(%d,%d)", newX, newY)
	}

	gc.bx = gc.x
	gc.by = gc.y
	gc.x = newX
	gc.y = newY

	return nil
}

func (gc *GopherController) GetCurrentCoordinate() (int, int) {
	return gc.x, gc.y
}

func (gc *GopherController) GetBeforeCoordinate() (int, int) {
	return gc.bx, gc.by
}
