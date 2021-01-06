package main

import (
	"errors"
	"fmt"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/spf13/cobra"
)

var descrambleCommand = &cobra.Command{
	Use:   "descramble",
	Short: "converts hskXX.asd files to usable zip files",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg.Debug = debug
		if len(args) == 0 {
			return errors.New("pass .asd or .ncs file paths as argument")
		}

		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return err
		}

		if len(args) > 1 && out != "" {
			return errors.New("--out can be specified only with one input file")
		}

		for _, ifn := range args {
			ofn := ifn + ".zip"
			if out != "" && len(args) == 1 {
				ofn = out
			}

			if err := pkg.DecodeFileAndSave(ifn, ofn); err != nil {
				return fmt.Errorf("descramble of %s failed : %s", ifn, err)
			}
		}
		return nil
	},
}
