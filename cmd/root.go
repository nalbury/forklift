/*
Copyright © 2021 Nick Albury <nickalbury@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cmd contains our rootCmd
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/spf13/cobra"
)

// Config stores our configuration for mirroring
type Config struct {
	Source  string
	Dest    string
	Verbose bool
	Platform string
}

var conf Config

// parsePlatform parses the platform flag in the format "os/arch[/variant]"
// Taken mostly from https://github.com/google/go-containerregistry/blob/9cf3ed4ac182c640e59a82343512e4c2cacb6d68/cmd/crane/cmd/util.go#L10
func parsePlatform(platform string) (*v1.Platform, error) {
	if platform == "all" {
		return nil, nil
	}
	p := &v1.Platform{}
	parts := strings.Split(platform, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("failed to parse platform '%s': expected format os/arch[/variant]", platform)
	}
	if len(parts) > 3 {
		return nil, fmt.Errorf("failed to parse platform '%s': too many slashes", platform)
	}

	p.OS = parts[0]
	p.Architecture = parts[1]
	if len(parts) > 2 {
		p.Variant = parts[2]
	}

	return p, nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "forklift [flags] <src> <dst>",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse args
		if len(args) < 2 {
			// if err := cmd.Usage(); err != nil {
			// 	return err
			// }
			return fmt.Errorf("Please provide a source and dest registry\n")
		}
		conf.Source = args[0]
		conf.Dest = args[1]
		// Set up logging
		logs.Progress.SetOutput(os.Stdout)
		if conf.Verbose {
			logs.Debug.SetOutput(os.Stderr)
		}

		// options is a list of crane options for our copy operation
		var options []crane.Option

		// Parse platform and append to the list of options
		platform, err := parsePlatform(conf.Platform)
		if err != nil {
			return err
		}
		options = append(options, crane.WithPlatform(platform))

		// Copy image
		if err := crane.Copy(conf.Source, conf.Dest, options...); err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.forklift.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&conf.Verbose, "verbose", "v", false, "Enables verbose logging" )
	rootCmd.PersistentFlags().StringVarP(&conf.Platform, "platform", "p", "all", "Specify a specific os, arch and variant in the format os/arch[/variant] e.g. linux/arm64/v8")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
