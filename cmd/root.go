package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/semihbkgr/fmtdump/internal/format"
	"github.com/semihbkgr/fmtdump/internal/parse"
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

	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	p := parse.NewParser(file, f)
	for data, err := p.Next(); err != io.EOF; data, err = p.Next() {
		if err != nil {
			return err
		}
		for _, d := range data {
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %x\n", d.Block.Name, d.Value)
		}
	}

	return nil
}
