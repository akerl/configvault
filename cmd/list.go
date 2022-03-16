package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/akerl/configvault/vault"
)

func listRunner(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	if len(args) != 2 {
		return fmt.Errorf("invalid args")
	}

	private, err := flags.GetBool("private")
	if err != nil {
		return err
	}

	query := vault.Query{
		Bucket: args[0],
		Key:    args[1],
		Public: !private,
	}
	result, err := vault.Search(query)
	if err != nil {
		return err
	}

	for _, host := range result {
		fmt.Println(host)
	}
	return nil
}

var listCmd = &cobra.Command{
	Use:   "list [BUCKET] [KEY]",
	Short: "List keys in ConfigVault",
	RunE:  listRunner,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("private", false, "List hosts with private data")
}
