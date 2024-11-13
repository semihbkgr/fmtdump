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
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("format", "f", "", "path of the format file")
}

func run(cmd *cobra.Command, args []string) error {
	formatFile, err := cmd.Flags().GetString("format")
	if err != nil {
		return err
	}
	f, err := format.ParseFormatFile(formatFile)
	if err != nil {
		return err
	}
	if err := f.Validate(); err != nil {
		return err
	}

	file, err := os.Open(cmd.Flags().Args()[0])
	if err != nil {
		return err
	}
	p := parse.NewParser(file, f)

	for data, err := p.Next(); err != io.EOF; data, err = p.Next() {
		if err != nil {
			return err
		}
		for _, d := range data {
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %x\n", d.Field.Name, d.Value)
		}
	}

	return nil
}
