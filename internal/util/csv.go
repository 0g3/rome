package util

import (
	"encoding/csv"
	"os"
	"strings"
)

func ConvertMaze(csvPath, dstPath string) error {
	// csv読み込む
	csvFile, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// mazeファイルとして書き込む
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, row := range rows[1 : len(rows)-2] {
		joined := strings.Join(row[1:], "") + "\n"
		if _, err := f.Write([]byte(joined)); err != nil {
			return err
		}
	}

	joined := strings.Join(rows[len(rows)-2][1:], "")
	if _, err = f.Write([]byte(joined)); err != nil {
		return err
	}

	return nil
}
