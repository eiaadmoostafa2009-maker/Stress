/*
Copyright © 2026 Eyad_Abdelazim eiaadmoostafa2009@gmail.com
*/

package cmd

import (
	"context"
	"fmt"
	"os"
	"stress/internal/config"
	"stress/internal/reporter"
	"stress/internal/service"
	"stress/internal/worker"

	"github.com/spf13/cobra"
)

// reqCmd represents the req command
var reqCmd = &cobra.Command{
	Use:   "req",
	Short: "Send a single request to test an endpoint",
	Long: `The command that allows you to test a specific endpoint 
	to test:-
	1-the service
	2-the route
	3-check errors`,
	RunE: func(cmd *cobra.Command, args []string) error{
		cfg := config.TestConfig{
			TargetUrl:   targetUrl,
			JWT: jwt,
	        Method: method,
	        Body: body,
		}

		worker := worker.NewWorker(cfg)
		reporter := reporter.NewReporter(os.Stdout)
		svc := service.NewService(nil, reporter, cfg, worker)
        
		if err:=svc.SendRequest(context.Background()); err != nil{
			return fmt.Errorf("run command: %w", err)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reqCmd)

	
	runCmd.Flags().StringVarP(
		&targetUrl, 
		"url", 
		"u", 
		"", 
		"The URL to test",
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
