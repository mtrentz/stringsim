/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
		// mainStrings containing the 's1's
		// otherStrings containing the 's2's, ...
		// All treated as a list since when using --f1 flag,
		// you can provide more than one s1.
		var mainStrings []string
		var otherStrings []string

		// FLAG LOGIC
		// First check if nothing was provided, if so
		// only print the usage message.
		if len(args) == 0 && File1 == "" && File2 == "" {
			cmd.Usage()
			return
		}
		// Quickly check if any input was provided its either
		// a txt file or a json extension
		if File1 != "" {
			if ext := filepath.Ext(File1); ext != ".json" && ext != ".txt" {
				fmt.Println("File1 extension not .json or .txt")
				os.Exit(1)
			}
		}
		if File2 != "" {
			if ext := filepath.Ext(File2); ext != ".json" && ext != ".txt" {
				fmt.Println("File2 extension not .json or .txt")
				os.Exit(1)
			}
		}
		// If File1 was provided, I either need at least
		// one argument (s2) or File2
		if File1 != "" {
			// Read 's1's from the file
			mainStrings = utils.ReadFromFile(File1)
			// Check if File2 was provided
			if File2 != "" {
				// Read 's2's from the file
				otherStrings = utils.ReadFromFile(File2)
			} else {
				// If File2 was not provided,
				// then only one argument (s2) needs to be provided
				utils.CheckForMinimumArgs(cmd, 1, mainStrings)
				// Read all 's2's from the arguments
				otherStrings = args
			}
		} else {
			// If File1 was not provided, I need to check for
			// File2
			if File2 != "" {
				// If File2 was provided, I need at least one argument
				// to be the s1.
				utils.CheckForMinimumArgs(cmd, 1, args)
				// I'll read 's2's from file
				otherStrings = utils.ReadFromFile(File2)
				// And 's1's from the arguments
				mainStrings = args
			} else {
				// If File1 and File2 were not provided,
				// I need at least two arguments
				utils.CheckForMinimumArgs(cmd, 2, args)
				// s1 will be the first
				mainStrings = []string{args[0]}
				// s2 will be the rest
				otherStrings = args[1:]
			}
		}
		// Check if output is either a .json or .csv
		if Output != "" {
			if ext := filepath.Ext(Output); ext != ".json" && ext != ".csv" {
				fmt.Println("Output file extension not .json or .csv")
				os.Exit(1)
			}
		}

		// Case insensitive and unidecode flags
		if Insensitive {
			utils.SliceToLower(&mainStrings)
			utils.SliceToLower(&otherStrings)
		}
		if Unidecode {
			utils.SliceToUnidecode(&mainStrings)
			utils.SliceToUnidecode(&otherStrings)
		}

		// metric logic
		var metric string
		// Decide on the metric to use if has a flag
		if metric != "" {
			// Make sure Metric is all lower case
			metric = strings.ToLower(metric)
		} else {
			metric = "jaro"
		}

		// The task will be done concurrently
		// where the amount of goroutines is the smaller of the
		// number of CPUs and the length of otherStrings
		MAX_CPU_CORES := runtime.NumCPU()
		amountGoroutines := utils.Min(len(otherStrings), MAX_CPU_CORES)

		// Now I'll take the otherStrings and split them into
		// 'amountGoroutines' slices, as evenly as possible.
		// The logic is that each goroutine will get one of these sub slices
		// and for each element will calculate the similarity
		// against the all the mainStrings.
		otherStringsSubSlices := utils.SliceSplit(otherStrings, amountGoroutines)

		amountComputations := len(mainStrings) * len(otherStrings)

		// Set a threshold for too many computations. If it's too high,
		// I'll have a separate flow, which will not hold too much
		// into memory and will be apending to the output file instead
		threshold := 100000
		tooManyComputations := amountComputations > threshold

		// Won't print to screen if too many computations
		if tooManyComputations && Output == "" {
			fmt.Printf("Too many similarities to comput and print to screen. Please use -o to output to file.\n")
			os.Exit(1)
		}

		// I'll pass the flags as a map to the "flows"
		// if anything needs to be accesed.
		stringFlags := map[string]string{
			"File1":  File1,
			"File2":  File2,
			"Output": Output,
			"Metric": metric,
		}
		boolFlags := map[string]bool{
			"Insensitive": Insensitive,
			"Silent":      Silent,
			"Unidecode":   Unidecode,
		}

		// Send them to the proper flow
		if !tooManyComputations {
			similarity.NormalFlow(mainStrings, otherStringsSubSlices, metric, amountGoroutines, stringFlags, boolFlags)
		} else {
			similarity.BigFileFlow(mainStrings, otherStringsSubSlices, metric, amountGoroutines, stringFlags, boolFlags)
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
var Unidecode bool
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
	rootCmd.Flags().BoolVarP(&Unidecode, "unidecode", "u", false, "If provided, will use unidecode to get ASCII transliterations of Unicode text")
}
