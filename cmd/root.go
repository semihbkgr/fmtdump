package cmd

import (
	"fmt"
	"os"

	"github.com/semihbkgr/fmtdump/internal/format"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fmtdump",
	Short: "formatted dump",
	Args:  cobra.MinimumNArgs(1),
	RunE:  run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("file", "f", "", "path of the file")
}

func run(cmd *cobra.Command, args []string) error {
	f, err := format.ParseFormat(cmd.Flags().Args())
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStderr(), "%s\n", f)
	return nil
}
