package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/semihbkgr/fmtdump/internal/format"
	"github.com/semihbkgr/fmtdump/internal/parse"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                   "fmtdump --format=<format.json> <data-file>",
	Short:                 "Flexible data file dump tool for custom formats",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	RunE:                  run,
	Version:               version(),
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

	for entry, err := p.Next(); err != io.EOF; entry, err = p.Next() {
		if err != nil {
			return err
		}
		s, err := entry.String()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), s)
	}

	return nil
}

func version() string {
	//todo
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return info.Main.Version
}
