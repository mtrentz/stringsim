/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mtrentz/similarity-cli/similarity"
	"github.com/mtrentz/similarity-cli/utils"
)

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

		// Passes the main string normally as the first paramenter
		utils.CheckForMinimumArgs(cmd, 2, args)
		mainStr = args[0]
		otherStrings = args[1:]

		var similarities []similarity.Similarity

		for _, str := range otherStrings {
			sim := similarity.GetSimilarity(mainStr, str)
			similarities = append(similarities, sim)
			sim.Result()
		}

		if Output != "" {
			utils.WriteToFile(Output, similarities)
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

var File1 string
var File2 string
var Output string

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.similarity-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringVarP(&File1, "f1", "", "", "Path to input file containing many s1, to be compared against all other s2. This can be a .txt file separated by newlines, or a JSON list of strings.")
	rootCmd.Flags().StringVarP(&File2, "f2", "", "", "Path to input file containing many s2, to be compared against s1, many s1 in case f1 was provided. This can be a .txt file separated by newlines, or a JSON list of strings.")
	rootCmd.Flags().StringVarP(&Output, "out", "o", "", "Path to output file. If not provided, output will be printed to stdout.")
}
