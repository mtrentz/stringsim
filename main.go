/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"
	"time"

	"github.com/mtrentz/stringsim/cmd"
)

func main() {
	start := time.Now()
	cmd.Execute()
	duration := time.Since(start)

	// Formatted string, such as "2h3m0.5s" or "4.503μs"
	fmt.Println(duration)
}
