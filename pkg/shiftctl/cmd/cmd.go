/*
Copyright 2018 The Elasticshift Authors.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

func NewDefaultCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "shiftctl",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	
	rootCmd.AddCommand(commandPlugin())

	return rootCmd
}
