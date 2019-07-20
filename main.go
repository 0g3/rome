package main

import (
	"container/list"
	"errors"
	"github.com/0g3/rome/internal/rome"
	"time"
)

const (
	mazeW       = 31
	mazeH       = 31
	updateDelay = time.Millisecond * 100
)

type coordinate struct {
	x int
	y int
}

type coordinates []*coordinate

func (cs coordinates) CheckContain(crd *coordinate) bool {
	for _, c := range cs {
		if *c == *crd {
			return true
		}
	}
	return false
}

func (c *coordinate) move(d rome.MoveDirection) {
	switch d {
	case rome.MoveDirectionUp:
		c.y++
	case rome.MoveDirectionRight:
		c.x++
	case rome.MoveDirectionDown:
		c.y--
	case rome.MoveDirectionLeft:
		c.x--
	}
}

type dfsMover struct {
	stack            *list.List
	movedCoordinates coordinates

	// 先頭にあるものほど優先的に探索される
	directions []rome.MoveDirection
}

func NewDFSMover(d []rome.MoveDirection) *dfsMover {
	return &dfsMover{
		stack:            list.New(),
		movedCoordinates: make(coordinates, 0),
		directions:       d,
	}
}

func (dfsm *dfsMover) next(m rome.Maze, gc *rome.GopherController) bool {
	for _, d := range dfsm.directions {
		// 通過済みには移動しない
		x, y := gc.GetCurrentCoordinate()
		crd := &coordinate{x: x, y: y}
		crd.move(d)
		if dfsm.movedCoordinates.CheckContain(crd) {
			continue
		}

		if err := gc.Move(m, d); err == nil {
			bx, by := gc.GetBeforeCoordinate()
			bc := &coordinate{x: bx, y: by}
			dfsm.stack.PushBack(bc)
			dfsm.movedCoordinates = append(dfsm.movedCoordinates, bc)
			return true
		}
	}
	return false
}

func (dfsm *dfsMover) prev(m rome.Maze, gc *rome.GopherController) error {
	// 戻り先の座標を取得
	if dfsm.stack.Back() == nil {
		return errors.New("スタックが空")
	}
	crd := dfsm.stack.Back().Value.(*coordinate)
	dfsm.stack.Remove(dfsm.stack.Back())

	for _, d := range dfsm.directions {
		x, y := gc.GetCurrentCoordinate()
		crd2 := &coordinate{x: x, y: y}
		crd2.move(d)
		if crd2.x == crd.x && crd2.y == crd.y {
			if err := gc.Move(m, d); err != nil {
				return err
			}
			bx, by := gc.GetBeforeCoordinate()
			bc := &coordinate{x: bx, y: by}
			dfsm.movedCoordinates = append(dfsm.movedCoordinates, bc)
		}
	}
	return nil
}

func (dfsm *dfsMover) Move(m rome.Maze, gc *rome.GopherController) error {
	if !dfsm.next(m, gc) {
		if err := dfsm.prev(m, gc); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	m := rome.NewMaze(mazeW, mazeH)
	if err := m.LoadMazeFile("maze/maze1"); err != nil {
		panic(err)
	}

	dfsm := NewDFSMover([]rome.MoveDirection{
		rome.MoveDirectionUp,
		rome.MoveDirectionRight,
		rome.MoveDirectionDown,
		rome.MoveDirectionLeft,
	})

	e, err := rome.NewEmulator(m, dfsm, updateDelay)
	if err != nil {
		panic(err)
	}

	if err := e.Run(); err != nil {
		panic(err)
	}
}
