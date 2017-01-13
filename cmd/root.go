package cmd

import (
	"fmt"
	"os"

	"github.com/dtan4/valecjr/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	defaultTableName = "valec"
)

var rootOpts = struct {
	debug     bool
	tableName string
}{}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	Use:           "valecjr",
	Short:         "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := aws.Initialize(); err != nil {
			return errors.Wrap(err, "Failed to initialize AWS API clients.")
		}

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if rootOpts.debug {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Println(err)
		}
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&rootOpts.debug, "debug", false, "Debug mode")
	RootCmd.PersistentFlags().StringVar(&rootOpts.tableName, "table-name", defaultTableName, "DynamoDB table name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
