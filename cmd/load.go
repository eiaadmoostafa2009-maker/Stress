/*
Copyright © 2026 Eyad_Abdelazim eiaadmoostafa2009@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"stress/internal/config"
	"stress/internal/engine"
	"stress/internal/reporter"
	"stress/internal/service"
	"stress/internal/worker"
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
	Use:   "load",
	Short: "the command that starts the whole test",
	Long: `the command that config the url, number of workers and duration
	to start the whole test and return the results`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.TestConfig{
			TargetUrl:   targetUrl,
			Concurrency: concurrency,
			Duration:    duration,
			JWT: jwt,
	        Method: method,
	        Body: body,
		}

		//logic
		worker := worker.NewWorker(cfg)
		engine := engine.NewEngine(cfg, worker)
		reporter := reporter.NewReporter(os.Stdout)
		svc := service.NewService(engine, reporter, cfg, worker)

		err := svc.Load(context.Background())

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
