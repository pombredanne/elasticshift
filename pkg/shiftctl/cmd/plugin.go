/*
Copyright 2018 The Elasticshift Authors.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func commandPlugin() *cobra.Command {
	pluginCmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manages the shift plugins",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			// fmt.Printf("Manages the plugins....")
		},
	}

	pluginCmd.AddCommand(initCmd())
	pluginCmd.AddCommand(searchCmd())
	pluginCmd.AddCommand(pushCmd())

	return pluginCmd
}

func initCmd() *cobra.Command {

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create a shift plugin",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\n Init shift plugin.")
		},
	}
	return initCmd
}

func searchCmd() *cobra.Command {

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Perform a search against shift plugin registry",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\n Search shift plugin.")
		},
	}
	return searchCmd
}

func pushCmd() *cobra.Command {

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push the plugin to shift registry",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\n Push shift plugin.")
		},
	}
	return pushCmd
}
