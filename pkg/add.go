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
	"sort"

	"github.com/markbates/pkger"
)

func update(filePath string, gamePaths []zipFilePath) error {
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

func Add(rootDir string, categoryID int, newGames []*GameItem, fontName string, bgName string) ([]*GameItem, error) {
	if rootDir == "" {
		return nil, errors.New("root path value is empty")
	}
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("root path '%s' not exists", rootDir)
	}
	gamesDirectoryName := fmt.Sprintf("Game%02d", categoryID)
	gamesPath := path.Join(rootDir, gamesDirectoryName)
	if _, err := os.Stat(gamesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path '%s' doesn't exist", gamesPath)
	}
	hsk00Path := path.Join(gamesPath, "hsk00.asd")
	b, err := getHsk00lstContent(hsk00Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read '%s': %s", hsk00Path, err)
	}
	gameList, err := ParseHsk00lstContent(b)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hsk00 content: %s", err)
	}

	if Debug {
		fmt.Println("current game list:")
		for _, game := range gameList {
			fmt.Printf("%03d. %#v\n", game.ID, game)
		}
	}

	initialEndIndex := 5 - (len(gameList) % 5)
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
		var batchFilePaths []zipFilePath
		for _, g := range newGames[startIndex:y] {
			filename := sanitizeFilename(g.SourcePath)
			g.Filename = filename
			batchFilePaths = append(batchFilePaths, zipFilePath{filename: filename, srcPath: g.SourcePath})
		}

		hskID, err := gameList.NextHskID()
		if err != nil {
			return nil, fmt.Errorf("failed to get next hsk id: %s", err)
		}
		hskFileName := fmt.Sprintf("Hsk%02d.asd", hskID)
		hskFilePath := path.Join(gamesPath, hskFileName)
		// update only last partial hsk files. create all other files from scratch
		if _, err := os.Stat(hskFilePath); err == nil && endIndex < 5 {
			if err := update(hskFilePath, batchFilePaths); err != nil {
				return nil, fmt.Errorf("failed to update games to file: %s: %s", hskFilePath, err)
			}
		} else {
			if err := GenerateScrambledZip(batchFilePaths, hskFilePath); err != nil {
				return nil, fmt.Errorf("failed to create file: %s: %s", hskFilePath, err)
			}
		}

		for _, game := range batch {
			game.Hsk = hskFileName
			gameList = append(gameList, game)
		}

		hskID++
	}

	binPrefixF, err := pkger.Open("/assets/binprefix")
	if err != nil {
		return nil, fmt.Errorf("failed to open binprefix file: %s", err)
	}
	binPrefix, err := ioutil.ReadAll(binPrefixF)
	if err != nil {
		return nil, fmt.Errorf("failed to read binprefix: %s", err)
	}

	sort.Sort(gameList)
	for i := 0; i < len(gameList); i += 10 {
		start := i
		end := (i + 10)
		pageNo := (i / 10) + 1
		if end > len(gameList) {
			end = len(gameList)
		}
		BGFilename := fmt.Sprintf("Game%02d.bin", pageNo)

		gameListInPage := gameList[start:end]
		menuItemTexts := []string{}
		for i, gameItem := range gameListInPage {
			gameItem.BGFilename = BGFilename
			gameItem.ID = start + i + 1
			if gameItem.Filename == "" {
				return nil, fmt.Errorf("filename not present")
			}
			menuItemTexts = append(menuItemTexts, fmt.Sprintf("%02d. %s", gameItem.ID, gameItem.Name))
		}

		imageBytes, err := generateMenuImage(menuItemTexts, fontName, bgName)
		if err != nil {
			return nil, fmt.Errorf("menu image generation failed: %s", err)
		}

		BGFilePath := path.Join(gamesPath, BGFilename)
		if Debug {
			ioutil.WriteFile(BGFilePath+".debug.jpg", imageBytes, 0644)
		}

		scrambledBytes := append(binPrefix, imageBytes[2:]...)
		if err := ioutil.WriteFile(BGFilePath, scrambledBytes, 0644); err != nil {
			return nil, nil
		}
	}

	if Debug {
		log.Println("new game list:")
		for _, g := range gameList {
			log.Printf("%#v\n", g)
		}
	}

	listZipBytes, err := makeHsk00(gameList)
	if err != nil {
		return nil, err
	}

	hsk00FilePath := path.Join(gamesPath, "hsk00.asd")
	if err := ioutil.WriteFile(hsk00FilePath, PKToWQW(listZipBytes), 0644); err != nil {
		return nil, err
	}

	return gameList, nil
}
