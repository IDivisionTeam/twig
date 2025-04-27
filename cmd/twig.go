package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"twig/config"
	"twig/log"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	twigCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		fmt.Sprintf(
			"config file (default is ~%s/%s/%s)",
			config.Path,
			config.Name,
			config.Type,
		),
	)

	twigCmd.AddCommand(
		createCmd,
		cleanCmd,
	)
}

func initConfig() {
	config.InitConfig(cfgFile)
}

var twigCmd = &cobra.Command{
	DisableAutoGenTag: true,
	Use:               "twig",
	Version:           "1.3.0",
	Args:              cobra.NoArgs,
}

func Execute() {
	if err := twigCmd.Execute(); err != nil {
		log.Fatal().Println(err)
	}
}
