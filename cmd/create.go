package cmd

import (
	"github.com/spf13/cobra"
	"net/http"
	"twig/command"
	"twig/common"
	"twig/log"
	"twig/network"
)

var branchType string

func init() {
	createCmd.Flags().StringVarP(
		&branchType,
		"type",
		"t",
		"",
		"(optional) overrides the type of branch",
	)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create branch from Jira Issue",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		httpClient := &http.Client{}
		client := network.NewClient(httpClient)

		input := &common.Input{
			Flags:     common.EmptyFlag,
			Arguments: make(map[common.InputType]string),
		}

		input.Arguments[common.Issue] = args[0]

		if branchType != "" {
			input.Arguments[common.BranchType] = branchType
		}
		createCommand := command.NewCreateLocalBranchCommand(client, input)

		if err := createCommand.Execute(); err != nil {
			log.Fatal().Println(err)
		}
	},
}
