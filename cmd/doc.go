/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	flower "github.com/alzabo/flower/pkg"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate docs for flows",
	Long:  `Generate documentation for flows`,
	Run: func(cmd *cobra.Command, args []string) {
		err := flower.FlowDocsFromDirectories(dirs)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(docCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
