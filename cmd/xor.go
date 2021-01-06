package main

import (
	"io/ioutil"
	"os"

	"github.com/dev-drprasad/hsk00/pkg"
	"github.com/spf13/cobra"
)

var xorCommand = &cobra.Command{
	Use:   "xor",
	Short: "Just some debug helper",
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg.Debug = debug
		in, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		b, err := ioutil.ReadFile(in)
		if err != nil {
			return err
		}
		for _, c := range b {
			r := c ^ 0xe5
			os.Stdout.Write([]byte{r})
		}
		return nil
	},
}
