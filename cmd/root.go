/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var dirs []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flower",
	Short: "Utility for working with Concord flows",
	Long: `flower
	
Work with Concord flows`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.flower.yaml)")
	rootCmd.PersistentFlags().StringSliceVarP(&dirs, "dir", "d", []string{"./"}, "Directories to search for flow files")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().StringSliceP("dir", "d", []string{"./"}, "Directories to search for flow files")
}
