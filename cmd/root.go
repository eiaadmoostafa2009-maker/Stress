/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stress",
	Short: "A load testing app to test your app under stress",
	Long: `A load testing app that is developed to put your app under 
	pressure and see if it can stand while been used by a specific number
	of users at the same time`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {

}
