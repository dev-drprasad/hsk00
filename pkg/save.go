package pkg

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func Save(rootDir string, categoryID int, games []*GameItem, fontName string, bgName string) ([]*GameItem, error) {
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

	deletes := map[string]GameItemList{}
	for _, g := range games {
		if g.Deleted {
			deletes[g.Hsk] = append(deletes[g.Hsk], g)
		}
	}

	if Debug {
		fmt.Printf("games to delete %#v\n", deletes)
	}

	for hsk, gameList := range deletes {
		if err := Remove(rootDir, categoryID, hsk, gameList); err != nil {
			// broken. restore with backup
			return nil, fmt.Errorf("failed to remove games from %s <- %s", hsk, err)
		}
	}

	var gameListAfterDelete GameItemList
	var newGames GameItemList
	for _, g := range games {
		if !g.Deleted {
			gameListAfterDelete = append(gameListAfterDelete, g)
		}
		if g.Hsk == "" {
			newGames = append(newGames, g)
		}
	}

	if Debug {
		fmt.Printf("games to add %#v\n", newGames)
	}

	nextHskID, err := gameListAfterDelete.NextHskID()
	if err != nil {
		return nil, fmt.Errorf("failed to get next hsk id <- %s", err)
	}
	for startIndex := 0; startIndex < len(newGames); startIndex += 5 {
		endIndex := startIndex + 5
		if endIndex > len(newGames) {
			endIndex = len(newGames)
		}

		batch := newGames[startIndex:endIndex]
		var batchFilePaths []zipFilePath
		for _, g := range batch {
			filename := sanitizeFilename(g.SourcePath)
			g.Filename = filename
			batchFilePaths = append(batchFilePaths, zipFilePath{filename: filename, srcPath: g.SourcePath})
		}

		hskFileName := fmt.Sprintf("Hsk%02d.asd", nextHskID)
		hskFilePath := path.Join(gamesPath, hskFileName)

		if err := GenerateScrambledZip(batchFilePaths, hskFilePath); err != nil {
			return nil, fmt.Errorf("failed to create file: %s: %s", hskFilePath, err)
		}

		for _, game := range batch {
			game.Hsk = hskFileName
			gameListAfterDelete = append(gameListAfterDelete, game)
		}
		nextHskID++
	}

	gameList, err := GenerateMenu(gameListAfterDelete, gamesPath, fontName, bgName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate menu <- %s", err)
	}

	if err := SaveGameList(gamesPath, gameList); err != nil {
		return nil, fmt.Errorf("failed to save game list <- %s", err)
	}

	return gameList, nil
}
