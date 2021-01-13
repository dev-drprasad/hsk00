package pkg

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var gamedirrx = regexp.MustCompile(`Game\d{2,}`)

func GetGameList(rootDir string, categoryID int) (map[string][]GameItem, error) {
	if rootDir == "" {
		return nil, errors.New("root path value is empty")
	}
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("root path '%s' not exists", rootDir)
	}

	hsk00Paths := map[string]string{}

	if categoryID == -1 {
		filepath.Walk(rootDir, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("path '%s' read failed: %s", path, err)
			}

			if !fi.IsDir() {
				return nil
			}
			if !gamedirrx.MatchString(fi.Name()) {
				return nil
			}

			hsk00Paths[fi.Name()] = filepath.Join(path, "hsk00.asd")
			return nil
		})
	} else {
		categoryDirName := fmt.Sprintf("Game%02d", categoryID)
		categoryDirPath := path.Join(rootDir, categoryDirName)
		if _, err := os.Stat(categoryDirPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("path '%s' doesn't exist", categoryDirPath)
		}

		hsk00Paths[categoryDirName] = path.Join(categoryDirPath, "hsk00.asd")
	}

	if len(hsk00Paths) == 0 {
		return nil, errors.New("weird, no hsk00 paths are given")
	}

	listMap := map[string][]GameItem{}

	for categoryDirName, hsk00Path := range hsk00Paths {
		b, err := getHsk00lstContent(hsk00Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read '%s': %s", hsk00Path, err)
		}
		list, err := ParseHsk00lstContent(b)
		if err != nil {
			return nil, fmt.Errorf("failed to parse hsk00 content: %s", err)
		}
		listMap[strings.Replace(categoryDirName, "Game", "Category ", 1)] = list
	}
	return listMap, nil
}
