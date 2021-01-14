package main

import (
	"fmt"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/spf13/cobra"
)

var gameListCommand = &cobra.Command{
	Use:   "list",
	Short: "list all games",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg.Debug = debug

		rootDir, err := cmd.Flags().GetString("root")
		if err != nil {
			return err
		}

		categoryID, err := cmd.Flags().GetInt("category")
		if err != nil {
			return err
		}

		if !cmd.Flags().Changed("category") {
			categoryID = -1
		}

		listMap, err := pkg.GetGameList(rootDir, categoryID)
		if err != nil {
			return fmt.Errorf("failed to get game list: %s", err)
		}

		total := 0
		for category, list := range listMap {
			fmt.Printf("========= %s =========\n", category)
			for _, i := range list {
				fmt.Printf("%02d %s", i.ID, i.Name)
				if debug {
					fmt.Printf(" %s (%s, %s)", i.Filename, i.Hsk, i.BGFilename)
				}
				fmt.Print("\n")
				total++
			}
			fmt.Print("\n")
		}
		if categoryID == -1 {
			fmt.Printf("total: %d\n", total)
		}
		return nil
	},
}

var gameCommand = &cobra.Command{
	Use:   "game",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
