package rome

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

type Tile int

const (
	TileEmpty Tile = iota
	TileRoad
	TileWall
	TileGoal
	tileGopher
)

type Maze [][]Tile

func NewMaze(width, height int) Maze {
	maze := make(Maze, height)
	for i := 0; i < height; i++ {
		maze[i] = make([]Tile, width)
	}
	return maze
}

func (m Maze) LoadMazeFile(path string) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	pad := 0
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			byt := buf[m.Width()*y+x+pad]
			if byt == 10 { // LF
				pad++
				byt = buf[m.Width()*y+x+pad]
			}

			var t Tile
			switch string(byt) {
			case strconv.Itoa(int(TileRoad)):
				t = TileRoad
			case strconv.Itoa(int(TileWall)):
				t = TileWall
			case strconv.Itoa(int(TileGoal)):
				t = TileGoal
			}

			if ok := m.Set(x, m.Height()-y-1, t); !ok {
				return fmt.Errorf("タイルのセットに失敗: (x,y)=(%d,%d)",
					x, m.Width()*(m.Height()-y-1))
			}
		}
	}

	return nil
}

func (m Maze) Set(x, y int, t Tile) bool {
	if !m.validateCoordinate(x, y) {
		return false
	}
	if m[y][x] == tileGopher {
		return false
	}
	m[y][x] = t
	return true
}

func (m Maze) setGopher(gc *GopherController) (bool, bool) {
	var goal bool

	if !m.validateCoordinate(gc.bx, gc.by) {
		return false, goal
	}
	m[gc.by][gc.bx] = TileRoad

	t, ok := m.Get(gc.x, gc.y)
	if !ok {
		return false, goal
	}
	if t == TileGoal {
		goal = true
	}

	if ok := m.Set(gc.x, gc.y, tileGopher); !ok {
		m[gc.by][gc.bx] = tileGopher
		return false, goal
	}

	return true, goal
}

func (m Maze) Get(x, y int) (Tile, bool) {
	if !m.validateCoordinate(x, y) {
		return TileEmpty, false
	}
	return m[y][x], true
}

func (m Maze) validateCoordinate(x, y int) bool {
	return x >= 0 && x < len(m[0]) && y >= 0 && y < len(m)
}

func (m Maze) Width() int {
	return len(m[0])
}

func (m Maze) Height() int {
	return len(m)
}
