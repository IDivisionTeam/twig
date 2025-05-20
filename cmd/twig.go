package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"twig/common"
	"twig/config"
	"twig/log"
)

const version = "1.4.3"

var (
	cfgFile string
	twigCmd = &cobra.Command{
		DisableAutoGenTag: true,
		Use:               "twig",
		Version:           version,
		Args:              cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			command, _, err := cmd.Find(os.Args[1:])
			if err != nil {
				return
			}

			cleanAllName := strings.HasPrefix(
				command.Name(),
				fmt.Sprintf("%s %s", cleanCmdName, cleanAllCmdName),
			)
			cleanLocalName := strings.HasPrefix(
				command.Name(),
				fmt.Sprintf("%s %s", cleanCmdName, cleanLocalCmdName),
			)

			createName := strings.HasPrefix(
				command.Name(),
				createCmdName,
			)

			matchesCmdName := cleanAllName || cleanLocalName || createName

			if command != nil && matchesCmdName {
				if !common.HasGit() {
					logCmdFatal(errors.New("first, Git must be installed! https://git-scm.com/downloads/mac"))
				}
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	
	log.CreateRecorders()

	twigCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		fmt.Sprintf(
			"config file (default is ~%s)",
			config.GetDefaultConfigPath(),
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

	if init != nil && !strings.HasPrefix(init.Name(), InitCmdName) {
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
