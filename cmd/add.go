package main

import (
	"errors"
	"fmt"
	"os"

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

		return pkg.Add(rootDir, categoryID, args, fontName)
	},
}
