package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"twig/log"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	twigCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is ~/.config/twig/twig.config)",
	)

	twigCmd.AddCommand(
		createCmd,
		cleanCmd,
	)
}

func initConfig() {
	if cfgFile != "" {
		log.Debug().Println("Using custom config")

		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("env")
	} else {
		log.Debug().Println("Using default config")

		home, err := homedir.Dir()
		if err != nil {
			log.Fatal().Println(err)
		}

		viper.AddConfigPath(home + "/.config/twig/")
		viper.SetConfigName("twig.config")
		viper.SetConfigType("env")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Println(err)
	}
	log.Debug().Println("Config loaded")
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
