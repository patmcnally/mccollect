package cmd

import "github.com/spf13/cobra"

var collectionName string

var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "Manage pack ownership in a collection",
}

func init() {
	collectionCmd.PersistentFlags().StringVar(&collectionName, "name", "default", "collection name")
	rootCmd.AddCommand(collectionCmd)
}
