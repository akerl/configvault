package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/akerl/configvault/vault"
)

func writeRunner(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	if len(args) != 2 && len(args) != 3 {
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

	var data string
	if len(args) == 3 {
		data = args[2]
	} else {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			return fmt.Errorf("no stdin provided")
		}
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		data = string(b)
	}

	path := vault.Path{
		Bucket: args[0],
		Key:    args[1],
		Public: public,
		User:   user,
	}
	return vault.Write(path, data)
}

var writeCmd = &cobra.Command{
	Use:   "write [BUCKET] [KEY] ([DATA])",
	Short: "Write key to ConfigVault",
	RunE:  writeRunner,
}

func init() {
	rootCmd.AddCommand(writeCmd)
	writeCmd.Flags().Bool("public", false, "Write to public data for host")
	writeCmd.Flags().StringP("user", "u", "", "User to write from")
}
