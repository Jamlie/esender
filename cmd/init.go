/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a configuration file for the email sender.",
	Long:  `Initializes a configuration file for the email sender in the home directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(path.Join(homedir, ".esender.toml"))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		file.WriteString("[smtp]\n")
		file.WriteString("email = \"\"\n")
		file.WriteString("password = \n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
