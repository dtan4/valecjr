package cmd

import (
	"fmt"
	"strings"

	"github.com/dtan4/valec/util"
	"github.com/dtan4/valecjr/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var dumpOpts = struct {
	dotenvTemplate string
	override       bool
	output         string
	quote          bool
}{}

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: doDump,
}

func doDump(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Please specify namespace.")
	}
	namespace := args[0]

	secrets, err := aws.DynamoDB.ListSecrets(rootOpts.tableName, namespace)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve secrets.")
	}

	if len(secrets) == 0 {
		return errors.Errorf("Namespace %s does not exist.", namespace)
	}

	var dotenv []string

	if dumpOpts.dotenvTemplate == "" {
		dotenv, err = dumpAll(secrets, dumpOpts.quote)
		if err != nil {
			return errors.Wrap(err, "Failed to dump all secrets.")
		}
	} else {
		dotenv, err = dumpWithTemplate(secrets, dumpOpts.quote, dumpOpts.dotenvTemplate, dumpOpts.override)
		if err != nil {
			return errors.Wrap(err, "Failed to dump secrets with dotenv template.")
		}
	}

	if dumpOpts.output == "" {
		for _, line := range dotenv {
			fmt.Println(line)
		}
	} else {
		body := []byte(strings.Join(dotenv, "\n") + "\n")
		if dumpOpts.override {
			if err := util.WriteFile(dumpOpts.output, body); err != nil {
				return errors.Wrapf(err, "Failed to write dotenv file. filename=%s", dumpOpts.output)
			}
		} else {
			if err := util.WriteFileWithoutSection(dumpOpts.output, body); err != nil {
				return errors.Wrapf(err, "Failed to write dotenv file. filename=%s", dumpOpts.output)
			}
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(dumpCmd)

	dumpCmd.Flags().BoolVar(&dumpOpts.override, "override", false, "Override values in existing template")
	dumpCmd.Flags().StringVarP(&dumpOpts.output, "output", "o", "", "File to flush dotenv")
	dumpCmd.Flags().BoolVarP(&dumpOpts.quote, "quote", "q", false, "Quote values")
	dumpCmd.Flags().StringVarP(&dumpOpts.dotenvTemplate, "template", "t", "", "Dotenv template")
}