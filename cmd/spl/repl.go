package main

import (
	"os"
	"os/user"

	"github.com/spf13/cobra"

	"github.com/lukasmalkmus/spl/internal/app/spl/repl"
)

// replCmd represents the repl command.
var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Read Evaluate Print Loop",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the current user.
		user, err := user.Current()
		if err != nil {
			return err
		}

		// Print some info and start the REPL.
		cmd.Printf("Hello %s! This is the Simple Programming Language.\n", user.Username)
		cmd.Printf("Feel free to type in commands.\n")
		return repl.Start(os.Stdin, cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}
