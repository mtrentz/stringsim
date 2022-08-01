/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mtrentz/stringsim/similarity"
	"github.com/mtrentz/stringsim/utils"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stringsim <s1> <s2> [<s3> ...] [flags]",
	Short: "Calculate the similarity between strings.",
	Long: `Calculate the similarity between at least two strings.

Comparing s1 to s2
  stringsim adam adan

Comparing s1 to s2 and s3, case insensitive, output result to file
  stringsim adam adan Aden -i -o output.json

Reading s2, s3, ..., from a txt file separated by newlines and comparing to 'adam' using Levenshtein as metric
  stringsim adam --f2 strings.txt -m Levenshtein

Reading many words from a json file (formated as array of strings ["a", "b", ...]) and comparing each to every word in a txt file separated by newlines.
  stringsim --f1 strings_one.json --f2 strings_two.txt
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// All treated as list, even if only one mainString
		var mainStrings []string
		var otherStrings []string

		// If File1 and File2 not provided, needed at least two arguments
		if File1 == "" && File2 == "" {
			utils.CheckForMinimumArgs(cmd, 2, args)
		}

		// If not File1, then mainString is the first argument
		if File1 == "" {
			mainStrings = append(mainStrings, args[0])
		}

		// If not File2, then otherStrings is the rest of the arguments
		if File2 == "" {
			otherStrings = append(otherStrings, args[1:]...)
		}

		// If f1 or f2 is provided, I have to read the list
		// of strings from the file
		if File1 != "" {
			mainStrings = utils.ReadFromFile(File1)
			// Check if longer than 1
			utils.CheckForMinimumArgs(cmd, 1, mainStrings)
		}
		if File2 != "" {
			otherStrings = utils.ReadFromFile(File2)
			// Check if longer than 1
			utils.CheckForMinimumArgs(cmd, 1, otherStrings)
		}

		amountComputations := len(mainStrings) * len(otherStrings)

		// Set a threshold for too many computations. If more than this,
		// i'll write to output without holding in memory
		// and will not print anything to stdout
		threshold := 100000
		tooManyComputations := amountComputations > threshold

		if tooManyComputations && Output == "" {
			fmt.Printf("Too many similarities to comput and print to screen. Please use -o to output to file.\n")
			os.Exit(1)
		}

		// Now depending on the amount of computations, I'll either
		// send it to the 'NormalFlow' or the 'BigFileFlow'.
		// For that, I'll first need a map with the flags to pass
		// to them.
		stringFlags := map[string]string{
			"File1":  File1,
			"File2":  File2,
			"Output": Output,
			"Metric": Metric,
		}
		boolFlags := map[string]bool{
			"Insensitive": Insensitive,
			"Silent":      Silent,
		}

		if !tooManyComputations {
			similarity.NormalFlow(mainStrings, otherStrings, stringFlags, boolFlags)
		} else {
			similarity.BigFileFlow(mainStrings, otherStrings, stringFlags, boolFlags)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var Insensitive bool
var File1 string
var File2 string
var Output string
var Metric string
var Silent bool

func init() {
	rootCmd.Flags().BoolVarP(&Insensitive, "insensitive", "i", false, "Use case insensitive comparison")
	rootCmd.Flags().StringVarP(&File1, "f1", "", "", "Path to input file containing many s1, to be compared against all other s2. This can be a .txt file separated by newlines, or a JSON list of strings")
	rootCmd.Flags().StringVarP(&File2, "f2", "", "", "Path to input file containing many s2, to be compared against s1, many s1 in case f1 was provided. This can be a .txt file separated by newlines, or a JSON list of strings")
	rootCmd.Flags().StringVarP(&Output, "out", "o", "", "Path to output file. If not provided, output will be printed to stdout")
	rootCmd.Flags().StringVarP(&Metric, "metric", "m", "", "Metric used to compare strings. Defaults to Jaro. Available: Jaro, Levenshtein, DamerauLevenshtein, Hamming")
	rootCmd.Flags().BoolVarP(&Silent, "silent", "s", false, "If provided, will not print the results to stdout")
}
