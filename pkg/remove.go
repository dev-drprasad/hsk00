package pkg

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func Remove(rootDir string, categoryID int, hsk string, gameList GameItemList) error {
	if rootDir == "" {
		return errors.New("root path value is empty")
	}
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return fmt.Errorf("root path '%s' not exists", rootDir)
	}
	gamesDirectoryName := fmt.Sprintf("Game%02d", categoryID)
	gamesPath := path.Join(rootDir, gamesDirectoryName)
	if _, err := os.Stat(gamesPath); os.IsNotExist(err) {
		return fmt.Errorf("path '%s' doesn't exist", gamesPath)
	}
	hskFilePath := filepath.Join(gamesPath, hsk)

	zipFileInfos, err := Descramble(hskFilePath)
	if err != nil {
		return fmt.Errorf("failed to decompress file %s: %s", hskFilePath, err)
	}

	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)
	zw := zip.NewWriter(outz)
	defer zw.Close()

main:
	for _, zipfi := range zipFileInfos {
		header, err := zip.FileInfoHeader(zipfi.FileInfo())
		if err != nil {
			return fmt.Errorf("failed to read header: %s", err)
		}

		for _, g := range gameList {
			if g.Filename == EncodeFileName(header.FileInfo().Name()) {
				continue main
			}
		}

		fr, _ := zipfi.Open()
		defer fr.Close()

		fw, err := zw.CreateHeader(header)
		if _, err = io.Copy(fw, fr); err != nil {
			return fmt.Errorf("failed to copy to writer: %s", err)
		}
	}

	if err := zw.Close(); err != nil {
		return err
	}

	zipBytes := bb.Bytes()
	scambledBytes := PKToWQW(zipBytes)
	return ioutil.WriteFile(hskFilePath, scambledBytes, 0644)
}
