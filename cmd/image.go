package main

import (
	"fmt"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/spf13/cobra"
)

var imageGetCommand = &cobra.Command{
	Use:   "get",
	Short: "get hidden image from .bin and .logXX files",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg.Debug = debug

		for _, containerFilePath := range args {
			_, err := pkg.GetImage(containerFilePath)
			if err != nil {
				return fmt.Errorf("failed to get game list: %s", err)
			}
		}

		return nil
	},
}

var imageCommand = &cobra.Command{
	Use:   "image",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
