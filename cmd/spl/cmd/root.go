package cmd

import (
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "spl",
	Short: "Simple programming language compiler toolchain",
	Long: `spl is a full-blown compiler toolchain for the simple programming
language described here:

https://homepages.thm.de/~hg52/lv/compiler/praktikum/SPL-1.2.html

> Documentation & Support: https://github.com/lukasmalkmus/spl
> Source & Copyright Information: https://github.com/lukasmalkmus/spl`,
}

// Execute sets up the root command and executes it. An error is returned but
// handled internaly by cobra. The calling function should not handle the error
// but fail gracefully.
func Execute() error {
	// Persistent flags available to the root command and all of its children.
	// Configuration flags which go hand in hand with the configuration
	// specified in the configuration file and environment. Only available to
	// the root command.
	rootCmd.PersistentFlags().String("config", "", "configuration file to use")
	rootCmd.PersistentFlags().Uint("format.indent", 4, "indentation used by the formatter")

	// Bind the configuration flags to viper expect for the config flag.
	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "config" {
			return
		}
		_ = viper.BindPFlag(flag.Name, flag)
	})

	// Bind matching environment variables to viper.
	viper.SetEnvPrefix("e46")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// On initialization, the configuration is loaded.
	cobra.OnInitialize(initConfig(rootCmd))

	// Silence the usage message of the root command.
	rootCmd.SilenceUsage = true

	return rootCmd.Execute()
}

// initConfig reads the configuration using the provider viper instance.
func initConfig(cmd *cobra.Command) func() {
	return func() {
		viper.SetConfigType("toml")

		// If a configuration file is explicitly specified use it. If not,
		// search for it in common locations.
		if configFile := viper.GetString("config"); configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			home, err := homedir.Dir()
			if err != nil {
				cmd.Println(errors.Wrap(err, "finding home directory"))
				os.Exit(1)
			}
			viper.SetConfigName("spl")
			viper.AddConfigPath(".")
			viper.AddConfigPath(home)
		}

		// Read in the configuration but ignore the error if the configuration
		// file was not found in the search locations defined above.
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				cmd.Println(errors.Wrap(err, "reading configuration file"))
				os.Exit(1)
			}
		}
	}
}
