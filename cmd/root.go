/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/antzucaro/matchr"
	"github.com/spf13/cobra"
)

// Check for a minimum amount of arguments, if not enough,
// prints help page, error and exits.
func checkForMinimumArgs(cmd *cobra.Command, n int, args []string) {
	if len(args) < n {
		cmd.Usage()
		fmt.Printf("Expected %d arguments, got %d.\n", n, len(args))
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "similarity-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var mainStr string
		var otherStrings []string

		if Target != "" {
			// Passes the main string separately from the flag
			checkForMinimumArgs(cmd, 1, args)

			mainStr = Target
			otherStrings = args
		} else {
			// Passes the main string normally as the first paramenter
			checkForMinimumArgs(cmd, 2, args)

			mainStr = args[0]
			otherStrings = args[1:]
		}

		for _, str := range otherStrings {
			fmt.Printf("Similarity between %s and %s: %.4f\n", mainStr, str, matchr.Jaro(mainStr, str))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var Target string

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.similarity-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringVarP(&Target, "target", "t", "", "Target string to compare against")
}
