package main

import (
	"errors"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/spf13/cobra"
)

var scrambleCommand = &cobra.Command{
	Use:   "scramble",
	Short: "Removed in favor of 'scrambled-zip'",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("'scramble' command is removed in favour of 'scrambled-zip'. Please use 'scrambled-zip' from now on")
	},
}

var scrambledZipCommand = &cobra.Command{
	Use:   "scrambled-zip",
	Short: "Creates scrambled zip file from given files",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg.Debug = debug
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return err
		}

		return pkg.GenerateScrambledZip(pkg.ZipFilePathsFromPaths(args), out)
	},
}
