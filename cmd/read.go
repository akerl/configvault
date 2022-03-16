package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/akerl/configvault/vault"
)

func readRunner(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	if len(args) != 2 {
		return fmt.Errorf("invalid args")
	}

	private, err := flags.GetBool("private")
	if err != nil {
		return err
	}

	user, err := flags.GetString("user")
	if err != nil {
		return err
	}

	result, err := vault.Read(vault.Path{
		Bucket: args[0],
		Key:    args[1],
		Public: !private,
		User:   user,
	})
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}

var readCmd = &cobra.Command{
	Use:   "read [BUCKET] [KEY]",
	Short: "Read key from ConfigVault",
	RunE:  readRunner,
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().Bool("private", false, "Read from private data for host")
	readCmd.Flags().StringP("user", "u", "", "User to read from")
}
