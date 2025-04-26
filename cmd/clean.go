package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"twig/command"
	"twig/common"
	"twig/log"
	"twig/network"
)

var assignee string

func init() {
	cleanCmd.Flags().StringVarP(
		&assignee,
		"assignee",
		"a",
		"",
		"(optional) provides assignee to delete branch",
	)

	cleanCmd.AddCommand(
		cleanLocalCmd,
		cleanAllCmd,
	)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Deletes branches which have Jira tickets in 'Done' state",
	Args:  cobra.NoArgs,
}

var cleanLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Deletes only local branches which have Jira tickets in 'Done' state",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		httpClient := &http.Client{}
		client := network.NewClient(httpClient)

		input := &common.Input{
			Flags:     common.EmptyFlag,
			Arguments: make(map[common.InputType]string),
		}

		if assignee != "" {
			input.Arguments[common.Assignee] = assignee
		}
		deleteCommand := command.NewDeleteLocalBranchCommand(client, input)

		if err := deleteCommand.Execute(); err != nil {
			log.Fatal().Println(err)
		}
	},
}

var cleanAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Deletes remote and local branches which have Jira tickets in 'Done' state",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		httpClient := &http.Client{}
		client := network.NewClient(httpClient)

		input := &common.Input{
			Flags:     common.EmptyFlag,
			Arguments: make(map[common.InputType]string),
		}

		input.Arguments[common.Remote] = viper.GetString("branch.origin")

		if assignee != "" {
			input.Arguments[common.Assignee] = assignee
		}
		deleteCommand := command.NewDeleteLocalBranchCommand(client, input)

		if err := deleteCommand.Execute(); err != nil {
			log.Fatal().Println(err)
		}
	},
}
