package cmd

import (
	"fmt"
	"os"

	"github.com/applandinc/appland-cli/internal/appland"
	"github.com/applandinc/appland-cli/internal/build"
	"github.com/applandinc/appland-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	api     *appland.Client
	rootCmd = &cobra.Command{
		Use:     "appland",
		Short:   "Manage AppLand resources",
		Version: build.Version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
)

func fail(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func Execute() {
	config.LoadCLIConfig()
	api = appland.MakeClient(config.GetCurrentContext())

	if err := rootCmd.Execute(); err != nil {
		fail(err)
	}

	if err := config.WriteCLIConfig(); err != nil {
		fail(err)
	}
}
