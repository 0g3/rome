package main

import (
	"github.com/0g3/rome/internal/util"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		panic("使い方ちゃう")
	}
	if err := util.ConvertMaze(os.Args[1], os.Args[2]); err != nil {
		panic(err)
	}
}
