package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"twig/config"
	"twig/log"
)

const version = "1.3.1"

var (
	cfgFile string
	twigCmd = &cobra.Command{
		DisableAutoGenTag: true,
		Use:               "twig",
		Version:           version,
		Args:              cobra.NoArgs,
	}
)

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

func Execute() {
	if err := twigCmd.Execute(); err != nil {
		log.Fatal().Println(err)
	}
}

func logCmdFatal(cmdName string, err error) {
	ew := fmt.Errorf("%s: %w", cmdName, err).Error()
	log.Fatal().Println(ew)
}
