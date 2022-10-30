/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/spf13/cobra"

	flower "github.com/alzabo/flower/pkg"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List flows",
	Long: `TODO: list flows. For example:

flower list -d ./somedir`,
	Run: func(cmd *cobra.Command, args []string) {
		flows, err := flower.FlowsFromDirectories(dirs)
		if err != nil {
			fmt.Println("An error occurred")
		}

		tpl, err := template.New("ls").Parse("{{ .Name }}\t{{ .Doc }}")
		if err != nil {
			panic(err)
		}

		for _, flow := range flows {
			buf := new(bytes.Buffer)
			tpl.Execute(buf, flow)
			fmt.Println(buf)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
