package cmd

import (
	"fmt"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configTemplate is the configuration file template.
const configTemplate = `# SPL COMPILER TOOLCHAIN CONFIGURATION

# Source code formatter configuration.
[format]
# Indentation width used.
indent = {{ .format.indent }}
`

// configCmd represents the config command.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the configuration",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if t, err := template.New("config").Parse(configTemplate); err != nil {
			return fmt.Errorf("invalid configuration template: %w", err)
		} else if err := t.Execute(cmd.OutOrStdout(), viper.AllSettings()); err != nil {
			return fmt.Errorf("execute configuration template: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
