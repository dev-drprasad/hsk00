package pkg

import (
	"io/ioutil"
	"log"
	"path"
)

func SaveGameList(gamesPath string, gameList GameItemList) error {
	if Debug {
		log.Println("new game list:")
		for _, g := range gameList {
			log.Printf("%#v\n", g)
		}
	}

	listZipBytes, err := makeHsk00(gameList)
	if err != nil {
		return err
	}

	hsk00FilePath := path.Join(gamesPath, "hsk00.asd")
	if err := ioutil.WriteFile(hsk00FilePath, PKToWQW(listZipBytes), 0644); err != nil {
		return err
	}

	return nil
}
