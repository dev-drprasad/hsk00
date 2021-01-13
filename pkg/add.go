package pkg

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dev-drprasad/hsk00/util"
	"github.com/markbates/pkger"
)

func update(filePath string, gamePaths []string) error {
	zipFileInfos, err := Descramble(filePath)
	if err != nil {
		return fmt.Errorf("failed to decompress file %s: %s", filePath, err)
	}

	bb := bytes.NewBuffer(nil)
	outz := bufio.NewWriter(bb)
	zw := zip.NewWriter(outz)
	defer zw.Close()

	for _, zipfi := range zipFileInfos {
		header, err := zip.FileInfoHeader(zipfi.FileInfo())
		if err != nil {
			return fmt.Errorf("failed to read header: %s", err)
		}

		fr, _ := zipfi.Open()
		defer fr.Close()

		fw, err := zw.CreateHeader(header)
		if _, err = io.Copy(fw, fr); err != nil {
			return fmt.Errorf("failed to copy to writer: %s", err)
		}
	}

	if err := CompressFilesAndWrite(zw, gamePaths); err != nil {
		return fmt.Errorf("failed to write to zip writer: %s", err)
	}

	if err := zw.Close(); err != nil {
		return err
	}

	zipBytes := bb.Bytes()
	scambledBytes := PKToWQW(zipBytes)
	return ioutil.WriteFile(filePath, scambledBytes, 0644)
}

func Add(rootDir string, categoryID int, newGames []string, fontName string) error {
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
	hsk00Path := path.Join(gamesPath, "hsk00.asd")
	menuList, err := getMenuList(hsk00Path)
	if err != nil {
		return fmt.Errorf("failed to get list from hsk00.asd: %s", err)
	}
	if Debug {
		log.Println("current game list:")
		for i, name := range menuList {
			log.Printf("%03d. %s\n", i+1, name)
		}
	}

	initialEndIndex := 5 - (len(menuList) % 5)
	for endIndex := initialEndIndex; endIndex <= len(newGames)+5; endIndex += 5 {
		y := endIndex
		// dont mutate endIndex, will cause ∞ loop
		if endIndex > len(newGames) {
			y = len(newGames)
		}
		startIndex := endIndex - 5
		if startIndex < 0 {
			startIndex = 0
		}
		batch := newGames[startIndex:y]

		hskID := (len(menuList) / 5) + 1
		hskFileName := fmt.Sprintf("Hsk%02d.asd", hskID)
		hskFilePath := path.Join(gamesPath, hskFileName)
		// update only last partial hsk files. create all other files from scratch
		if _, err := os.Stat(hskFilePath); err == nil && endIndex < 5 {
			if err := update(hskFilePath, batch); err != nil {
				return fmt.Errorf("failed to update games to file: %s: %s", hskFilePath, err)
			}
		} else {
			if err := GenerateScrambledZip(batch, hskFilePath); err != nil {
				return fmt.Errorf("failed to create file: %s: %s", hskFilePath, err)
			}
		}

		newMenuImageFileNames := make(map[string]int) // simple set
		for _, gamePath := range batch {
			pageNo := len(menuList)/10 + 1
			menuImageFileName := fmt.Sprintf("Game%02d.bin", pageNo)
			menuListItem := fmt.Sprintf("%s,%s,0,1,%s", hskFileName, filepath.Base(gamePath), menuImageFileName)
			menuList = append(menuList, menuListItem)
			newMenuImageFileNames[menuImageFileName] = pageNo
		}

		binPrefixF, err := pkger.Open("/assets/binprefix")
		if err != nil {
			return fmt.Errorf("failed to open binprefix file: %s", err)
		}
		binPrefix, err := ioutil.ReadAll(binPrefixF)
		if err != nil {
			return fmt.Errorf("failed to read binprefix: %s", err)
		}
		for fname, pageNo := range newMenuImageFileNames {
			start := (pageNo - 1) * 10
			end := pageNo * 10
			if end > len(menuList) {
				end = len(menuList)
			}
			menuListInPage := menuList[start:end]
			gameNames := []string{}
			for i, menuItem := range menuListInPage {
				gameFilename := strings.SplitN(menuItem, ",", 3)[1]
				gameFileName := util.FileNameWithoutExt(gameFilename)
				gameName := strings.ReplaceAll(gameFileName, "_", " ")
				gameNames = append(gameNames, fmt.Sprintf("%02d.%s", start+i+1, gameName))
			}

			imageBytes, err := generateMenuImage(gameNames, fontName)
			if err != nil {
				return fmt.Errorf("menu image generation failed: %s", err)
			}

			filePath := path.Join(gamesPath, fname)
			if Debug {
				ioutil.WriteFile(filePath+".debug.jpg", imageBytes, 0644)
			}

			scrambledBytes := append(binPrefix, imageBytes[2:]...)
			if err := ioutil.WriteFile(filePath, scrambledBytes, 0644); err != nil {
				return nil
			}
		}
		hskID++
	}

	if Debug {
		log.Println("new game list:")
		log.Printf("%s\n", strings.Join(menuList, "\n"))
	}

	listZipBytes, err := makeHsk00(menuList)
	if err != nil {
		return err
	}

	hsk00FilePath := path.Join(gamesPath, "hsk00.asd")
	if err := ioutil.WriteFile(hsk00FilePath, PKToWQW(listZipBytes), 0644); err != nil {
		return err
	}

	return nil
}
