package main

import (
	"errors"
	"fmt"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/spf13/cobra"
)

var descrambleCommand = &cobra.Command{
	Use:   "descramble",
	Short: "can extract hidden content from scrambled files. currently can extract zip files and images",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg.Debug = debug
		if len(args) == 0 {
			return errors.New("pass .asd|.ncs|.bin file paths as argument")
		}

		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return err
		}

		if len(args) > 1 && out != "" {
			return errors.New("--out can be specified only with one input file")
		}

		for _, ifn := range args {
			if err := pkg.DecodeFileAndSave(ifn, out); err != nil {
				fmt.Println(err.Error())
			}
		}
		return nil
	},
}
