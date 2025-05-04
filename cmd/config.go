package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"twig/config"
	"twig/log"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "You can query/set/replace options with this command. The name is the section and the key separated by a dot",
		Args:  cobra.NoArgs,
	}
	configListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all variables set in config file, along with their values",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.GetConfigSnapshot()
			if err != nil {
				logCmdFatal(err)
			}

			printString(config.ProjectHost, cfg.Project.Host)
			printString(config.ProjectAuth, cfg.Project.Auth)
			printString(config.ProjectEmail, cfg.Project.Email)
			printString(config.ProjectToken, cfg.Project.Token)

			printString(config.BranchDefault, cfg.Branch.Default)
			printString(config.BranchOrigin, cfg.Branch.Origin)
			printStringArr(config.BranchExclude, cfg.Branch.Exclude)

			printStringArr(config.MappingBuild, cfg.Mapping.Build)
			printStringArr(config.MappingChore, cfg.Mapping.Chore)
			printStringArr(config.MappingCi, cfg.Mapping.Ci)
			printStringArr(config.MappingDocs, cfg.Mapping.Docs)
			printStringArr(config.MappingFeat, cfg.Mapping.Feat)
			printStringArr(config.MappingFix, cfg.Mapping.Fix)
			printStringArr(config.MappingPref, cfg.Mapping.Pref)
			printStringArr(config.MappingRefactor, cfg.Mapping.Refactor)
			printStringArr(config.MappingRevert, cfg.Mapping.Revert)
			printStringArr(config.MappingStyle, cfg.Mapping.Style)
			printStringArr(config.MappingTest, cfg.Mapping.Test)
		},
	}
	configGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Emits the value of the specified key",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]
			token, err := config.FromInput(input)
			if err != nil {
				logCmdFatal(err)
			}

			if strings.Contains(input, config.FromToken(config.Mapping)) || input == config.FromToken(config.BranchExclude) {
				printStringArr(token, config.GetStringArray(token))
			} else {
				printString(token, config.GetString(token))
			}
		},
	}
	configSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set value for config option",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			value := args[1]
			token, err := config.FromInput(name)
			if err != nil {
				logCmdFatal(err)
			}

			if strings.Contains(name, config.FromToken(config.Mapping)) || name == config.FromToken(config.BranchExclude) {
				if err := config.SetStringArray(token, strings.Split(value, ",")); err != nil {
					logCmdFatal(err)
				}
			} else {
				if err := config.SetString(token, value); err != nil {
					logCmdFatal(err)
				}
			}
		},
	}
)

func init() {
	configCmd.AddCommand(
		configListCmd,
		configGetCmd,
		configSetCmd,
	)
}

func printString(token config.Token, value string) {
	log.Info().Print(fmt.Sprintf("%s=%s", config.FromToken(token), value))
}
func printStringArr(token config.Token, values []string) {
	log.Info().Print(fmt.Sprintf("%s=%s", config.FromToken(token), strings.Join(values, ",")))
}
