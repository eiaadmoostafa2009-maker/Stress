/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"stress/internal/script"

	"github.com/spf13/cobra"
)

// scriptCmd represents the script command
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "test your api with a configured script",
	Long: `Allow stress to test your api with a scrpt you have written already like this:
    
config = {
    name = "Local API journey",
    base_url = "http://127.0.0.1:8080",

    users = 5,
    duration_seconds = 15,

    timeout_seconds = 3,
    think_time_seconds = 1,
}

function scenario(user)
    local health = user:get("/health")

    assert(
        health.status == 200,
        "Expected /health to return HTTP 200"
    )

    local posts = user:get("/api/posts")

    assert(
        posts.status == 200,
        "Expected /api/posts to return HTTP 200"
    )
end

`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := script.Run(cmd.Context(), args[0]); err != nil {
			return fmt.Errorf("Script error: %w", err)
		}
		return nil

	},
}

func init() {
	rootCmd.AddCommand(scriptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scriptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scriptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
