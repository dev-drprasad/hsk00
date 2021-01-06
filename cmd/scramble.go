package main

import (
	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/spf13/cobra"
)

var scrambleCommand = &cobra.Command{
	Use:   "scramble",
	Short: "Creates scrambled hsk.asd file from given files",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg.Debug = debug
		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return err
		}

		return pkg.GenerateScrambledZip(args, out)
	},
}
