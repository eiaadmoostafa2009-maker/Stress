/*
Copyright © 2026 Eyad_Abdelazim eiaadmoostafa2009@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"stress/internal/config"
	"stress/internal/loadtest"
	"stress/internal/service"
	"time"

	"github.com/spf13/cobra"
)

var (
	targetUrl   string
	concurrency uint
	duration    time.Duration
	method string
	body string 
	jwt string
)
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "the command that starts the whole test",
	Long: `the command that config the url, number of workers and duration
	to start the whole test and return the results`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.TestConfig{
			TargetUrl:   targetUrl,
			Concurrency: concurrency,
			Duration:    duration,
		}

		//logic
		loadtest := loadtest.NewEngine()
		svc := service.NewLoadTestService(*loadtest, os.Stdout)

		err := svc.Run(context.Background(), cfg)

		if err != nil {
			return fmt.Errorf("run command: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(
		&targetUrl, 
		"url", 
		"u", 
		"", 
		"The URL to test",
	)
	runCmd.Flags().UintVarP(
		&concurrency, 
		"concurrency", 
		"c", 
		1, 
		"Number of concurrent requests",
	)
	runCmd.Flags().DurationVarP(
		&duration, 
		"duration", 
		"d", 30*time.Second, 
		"Duration of the test in seconds",
	)
	runCmd.Flags().StringVarP(
		&method, 
		"method", 
		"m", 
		"GET", 
		"request method: GET, POST, PUT, PATCH, DELETE",
	)
	runCmd.Flags().StringVarP(
        &body,
        "body",
		"b",
        "",
        "Raw request body",
    ) 
	runCmd.Flags().StringVarP(
        &jwt,
        "jwt",
		"j",
        "",
        "JWT access token for authenticated requests",
    )

	_ = runCmd.MarkFlagRequired("url")
}
