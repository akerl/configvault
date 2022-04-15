package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/akerl/configvault/vault"
)

func deleteRunner(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	if len(args) != 2 {
		return fmt.Errorf("invalid args")
	}

	public, err := flags.GetBool("public")
	if err != nil {
		return err
	}

	user, err := flags.GetString("user")
	if err != nil {
		return err
	}

	path := vault.Path{
		Bucket: args[0],
		Key:    args[1],
		Public: public,
		User:   user,
	}
	return vault.Delete(path)
}

var deleteCmd = &cobra.Command{
	Use:   "delete [BUCKET] [KEY]",
	Short: "Delete key from ConfigVault",
	RunE:  deleteRunner,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().Bool("public", false, "Delete from public data for host")
	deleteCmd.Flags().StringP("user", "u", "", "User to delete from")
}
