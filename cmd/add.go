package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "add game(s) to category",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg.Debug = debug
		if len(args) == 0 {
			return errors.New("pass games path as argument")
		}

		categoryID, err := cmd.Flags().GetInt("category")
		if err != nil {
			return err
		}

		if categoryID < 0 {
			return errors.New("category number must be 0 or greater")
		}

		rootDir, err := cmd.Flags().GetString("root")
		if err != nil {
			return err
		}
		if _, err := os.Stat(rootDir); os.IsNotExist(err) {
			return fmt.Errorf("root directory %s not exists", rootDir)
		}

		fontName, err := cmd.Flags().GetString("font")
		if err != nil {
			return err
		}

		bgName, err := cmd.Flags().GetString("background")
		if err != nil {
			return err
		}

		var gameList pkg.GameItemList
		for _, gamePath := range args {
			gameList = append(gameList, &pkg.GameItem{SourcePath: gamePath, Name: pkg.GameNameFromFilename(filepath.Base(gamePath))})
		}

		_, err = pkg.Add(rootDir, categoryID, gameList, fontName, bgName)
		return err
	},
}
