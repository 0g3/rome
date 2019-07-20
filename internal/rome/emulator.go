package rome

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"os"
	"path"
	"time"
)

var (
	roadImg   *ebiten.Image
	wallImg   *ebiten.Image
	goalImg   *ebiten.Image
	gopherImg *ebiten.Image
)

const (
	gopherImagePath = "./gopher.png"
	tileSize        = 16
)

func init() {
	dir := path.Dir(os.Args[0])
	fmt.Println(dir)
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}

	var err error

	// road
	roadImg, err = ebiten.NewImage(tileSize, tileSize, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	if err = roadImg.Fill(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}); err != nil {
		panic(err)
	}

	// wall
	wallImg, err = ebiten.NewImage(tileSize, tileSize, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	if err = wallImg.Fill(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}); err != nil {
		panic(err)
	}

	// goal
	goalImg, err = ebiten.NewImage(tileSize, tileSize, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	if err = goalImg.Fill(color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}); err != nil {
		panic(err)
	}

	gopherImg, _, err = ebitenutil.NewImageFromFile(gopherImagePath, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
}

type GopherMover interface {
	Move(Maze, *GopherController) error
}

type Emulator struct {
	m           Maze
	gc          *GopherController
	gm          GopherMover
	updateDelay time.Duration
	goal        bool
	moveCount   int
}

func NewEmulator(m Maze, gm GopherMover, updateDelay time.Duration) (*Emulator, error) {
	return &Emulator{
		m:           m,
		gc:          new(GopherController),
		gm:          gm,
		updateDelay: updateDelay,
	}, nil
}

func (e *Emulator) update(screen *ebiten.Image) error {
	if e.goal {
		if err := ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf("MOVE COUNT: %d", e.moveCount),
		); err != nil {
			return err
		}
		return nil
	}

	// screen 初期化
	if err := screen.Fill(color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}); err != nil {
		return err
	}

	// gopher setting
	ok, goal := e.m.setGopher(e.gc)
	if !ok {
		return errors.New("ゴーファーセッティングエラー")
	}
	e.goal = goal

	// map drawing
	for y := 0; y < e.m.Height(); y++ {
		for x := 0; x < e.m.Width(); x++ {
			t, ok := e.m.Get(x, y)
			if !ok {
				return errors.New("マップ描画エラー")
			}

			var img *ebiten.Image
			switch t {
			case TileEmpty:
				return fmt.Errorf("タイルがセットされていない部分が存在します: (x, y)=(%d, %d)", x, y)
			case TileRoad:
				img = roadImg
			case TileWall:
				img = wallImg
			case TileGoal:
				img = goalImg
			case tileGopher:
				img = gopherImg
			}
			if err := e.draw(screen, img, x, y); err != nil {
				return err
			}
		}
	}

	if err := e.gm.Move(e.m, e.gc); err != nil {
		return err
	}
	e.moveCount++

	time.Sleep(e.updateDelay)

	return nil
}

func (e *Emulator) Run() error {
	return ebiten.Run(e.update, e.m.Width()*tileSize, e.m.Height()*tileSize, 1.5, "rome")
}

func (e *Emulator) getOption(x, y int) *ebiten.DrawImageOptions {
	op := new(ebiten.DrawImageOptions)
	op.GeoM.Translate(float64(x)*float64(tileSize), float64(e.m.Height()-y-1)*float64(tileSize))
	return op
}

func (e *Emulator) draw(screen, i *ebiten.Image, x, y int) error {
	return screen.DrawImage(i, e.getOption(x, y))
}
