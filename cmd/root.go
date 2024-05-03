/*
Copyright © 2024 buty4649

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"time"

	"ndpp/ndp"

	"github.com/spf13/cobra"
)

type Options struct {
	timeout  int
	template string
	output   string
	command  string
}

var opts Options

var rootCmd = &cobra.Command{
	Use:          "ndpp [flags] ifname...",
	SilenceUsage: true,
	Short:        "A CLI tool for automating tasks based on Neighbor Discovery Protocol events ",
	Long:         "A CLI tool for automating tasks based on Neighbor Discovery Protocol events ",
	PreRunE:      validateOptions,
	RunE: func(cmd *cobra.Command, args []string) error {
		var result []*ndp.Result

		for _, name := range args {

			var ctx context.Context
			var cancel context.CancelFunc
			if opts.timeout == 0 {
				ctx, cancel = context.WithCancel(context.Background())
			} else {
				ctx, cancel = context.WithTimeout(context.Background(), time.Duration(opts.timeout)*time.Second)
			}
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			go func() {
				select {
				case <-sigChan:
					cancel()
				case <-ctx.Done():
					cancel()
				}
			}()

			r, err := ndp.SendRS(ctx, name)
			if err != nil {
				if err == context.DeadlineExceeded {
					return fmt.Errorf("timeout: %s", name)
				}
				return err
			}

			result = append(result, r)
		}

		if opts.command == "" && opts.template != "" {
			err := outputResults(result, opts.template, opts.output)
			if err != nil {
				return nil
			}
		}

		if opts.command != "" {
			err := commandExec(result, opts.command)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func validateOptions(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		cmd.SilenceErrors = true
		return fmt.Errorf("ifname is required")
	}

	if opts.timeout < 0 {
		return fmt.Errorf("timeout must be greater than or equal to 0")
	}

	if !slices.Contains([]string{"yaml", "json", "bird", "shellvar"}, opts.template) {
		return fmt.Errorf("template must be one of yaml, json, bird")
	}

	return nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&opts.timeout, "timeout", "t", 5, "timeout in seconds. if 0, it will wait indefinitely.")
	rootCmd.Flags().StringVarP(&opts.template, "template", "T", "yaml", "template possible values: yaml, json, bird, shellvar. default: yaml")
	rootCmd.Flags().StringVarP(&opts.output, "output", "o", "", "output file path. if not specified, it will be printed to stdout")
	rootCmd.Flags().StringVarP(&opts.command, "command", "c", "", "command to run when the router advertisement is received")

	rootCmd.MarkFlagsMutuallyExclusive("template", "command")
	rootCmd.MarkFlagsMutuallyExclusive("output", "command")
}
