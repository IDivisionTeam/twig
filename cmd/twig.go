package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"twig/config"
	"twig/log"
)

const version = "1.4.0"

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
		initCmd,
		createCmd,
		cleanCmd,
		configCmd,
	)
}

func initConfig() {
	init, _, err := twigCmd.Find(os.Args[1:])
	if err != nil {
		log.Debug().Println("Unable to find command for config initConfig")
		return
	}

	if init != nil && init.Name() != "init" {
		config.InitConfig(cfgFile)
	}
}

func Execute() {
	if err := twigCmd.Execute(); err != nil {
		logCmdFatal(err)
	}
}

func logCmdFatal(err error) {
	ew := fmt.Errorf("Error: %w", err)
	log.Fatal().Println(ew)
}
