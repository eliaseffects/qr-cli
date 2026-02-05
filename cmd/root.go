package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootFlags   OutputFlags
	showVersion bool

	rootCmd = &cobra.Command{
		Use:   "qr [data]",
		Short: "Generate QR codes from the terminal",
		Long: `qr-cli is a fast, minimal CLI for generating QR codes.

Examples:
  qr "https://example.com"              # Output to qr.png
  qr "Hello world" -o hello.png         # Custom output path
  echo "secret" | qr -o secret.png      # Read from stdin
  qr "https://example.com" --terminal   # Render in terminal
  qr "https://example.com" --open       # Open in viewer`,
		Args: cobra.MaximumNArgs(1),
		RunE: runRoot,
	}
)

func init() {
	addOutputFlags(rootCmd, &rootFlags, true)
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print version")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file path")
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if configErr != nil {
			return configErr
		}
		return nil
	}

	bindOutputFlags(rootCmd)

	rootCmd.AddCommand(wifiCmd)
	rootCmd.AddCommand(vcardCmd)
	rootCmd.AddCommand(batchCmd)
	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(versionCmd)
}

func runRoot(cmd *cobra.Command, args []string) error {
	applyOutputConfig(cmd, &rootFlags)

	if showVersion {
		printVersion(cmd)
		return nil
	}

	data := ""
	if len(args) > 0 {
		data = args[0]
	} else {
		input, err := readStdin()
		if err != nil {
			return err
		}
		data = input
	}

	return runGenerate(data, rootFlags, cmd.Flags().Changed("format"))
}

func readStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}
	return strings.TrimSpace(string(input)), nil
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
