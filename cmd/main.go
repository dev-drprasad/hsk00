package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hsk00",
	Short: "A tool to add/replace games to datafrog handheld console",
	Long:  `🚧 WIP 🚧`,
}

var debug bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug")

	xorCommand.Flags().String("file", "", "file to xor")
	rootCmd.AddCommand(xorCommand)

	descrambleCommand.Flags().String("out", "", "output zip file name (optional)")
	rootCmd.AddCommand(descrambleCommand)

	scrambledZipCommand.Flags().String("out", "", "output asd file name")
	scrambledZipCommand.MarkFlagRequired("out")
	rootCmd.AddCommand(scrambledZipCommand)

	// just need to tell people that we removed 'scramble' command
	scrambleCommand.Flags().String("out", "", "")
	rootCmd.AddCommand(scrambleCommand)

	addCommand.Flags().Int("category", 0, "number of category starting from 0, left -> right")
	addCommand.MarkFlagRequired("category")
	addCommand.Flags().String("root", "", "root path of sd card")
	addCommand.MarkFlagRequired("root")
	addCommand.Flags().String("font", "Gotham-Medium", "font name (Gotham-Medium | Video-Phreak) of menu text, default is Gotham-Medium")
	rootCmd.AddCommand(addCommand)

	gameCommand.PersistentFlags().String("root", "", "root path of sd card")
	gameCommand.MarkPersistentFlagRequired("root")
	gameListCommand.Flags().Int("category", 0, "number of category starting from 0, left -> right")
	gameCommand.AddCommand(gameListCommand)
	rootCmd.AddCommand(gameCommand)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
