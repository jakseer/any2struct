package main

import (
	"log"

	"github.com/jakseer/any2struct/cmd/any2struct/internal/generate"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "any2struct",
	Short:   "any2struct: convert any struct data to Go Struct Code",
	Long:    "any2struct: convert any struct data to Go Struct Code",
	Version: version,
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(generate.CmdGenerate)
}
