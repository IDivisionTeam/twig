package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"twig/branch"
	"twig/common"
	"twig/config"
	"twig/log"
	"twig/network"
)

var (
	branchType    string
	shouldPush    bool
	createCmdName = "create"
	createCmd     = &cobra.Command{
		Use:   createCmdName,
		Short: "Create branch from Jira Issue",
		Args:  cobra.MinimumNArgs(1),
		Run:   runCreate,
	}
)

func init() {
	createCmd.Flags().StringVarP(
		&branchType,
		"type",
		"t",
		"",
		"(optional) overrides the type of branch",
	)
	createCmd.Flags().BoolVarP(
		&shouldPush,
		"push",
		"p",
		false,
		"(optional) push branch to the remote",
	)
}

func runCreate(cmd *cobra.Command, args []string) {
	log.Debug().Println("create: executing command")

	httpClient := &http.Client{}
	client := network.NewHttpClient(httpClient)
	api:= network.NewJiraApi(client)
	issue := args[0]

	if err := validateIssue(issue); err != nil {
		logCmdFatal(err)
	}

	jiraIssue, err := api.GetJiraIssue(issue)
	if err != nil {
		logCmdFatal(err)
	}

	if err = validateBranchType(); err != nil {
		logCmdFatal(err)
	}

	excludePhrases := config.GetStringArray(config.BranchExclude)
	if len(excludePhrases) == 0 {
		log.Warn().Println(fmt.Sprintf("%q is not set", config.FromToken(config.BranchExclude)))
	}

	bt, err := convertInputToBranchType()
	if err != nil && branchType != "" {
		logCmdFatal(err)
	}

	b := branch.New(bt, excludePhrases)

	if b.Type == branch.NULL {
		jiraIssueTypes, err := api.GetJiraIssueTypes()
		if err != nil {
			logCmdFatal(err)
		}

		bt, err = convertIssueTypeToBranchType(*jiraIssue.Fields.Type, jiraIssueTypes)
		if err != nil {
			logCmdFatal(err)
		}

		b.Type = bt
	}

	branchName := b.BuildName(*jiraIssue)
	hasBranch := common.HasBranch(branchName)

	checkoutCommand, err := common.Checkout(branchName, hasBranch)
	if err != nil {
		logCmdFatal(err)
	}

	log.Info().Println(checkoutCommand)

	if shouldPush {
		remote := config.GetString(config.BranchOrigin)

		pushCommand, err := common.PushToRemote(branchName, remote)
		if err != nil {
			logCmdFatal(err)
		}

		log.Info().Println(pushCommand)
	}
}

func validateIssue(issue string) error {
	if issue == "" {
		return errors.New("validate: issue-key must not be empty")
	}

	return nil
}

func validateBranchType() error {
	log.Debug().Printf("create: validating type=%s", branchType)

	if branchType == "" {
		log.Debug().Println("create: type is empty, taking types from Jira")
		return nil
	}

	_, err := branch.InputToBranchType(branchType)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	log.Debug().Printf("create: type %s is valid", branchType)
	return nil
}

func convertInputToBranchType() (branch.Type, error) {
	bt, err := branch.InputToBranchType(branchType)
	return bt, fmt.Errorf("convert: %w", err)
}

func convertIssueTypeToBranchType(jiraIssueType network.IssueType, networkTypes []network.IssueType) (branch.Type, error) {
	mappedIssueTypes, err := branch.ConvertIssueTypesToMap(networkTypes)
	if err != nil {
		return branch.NULL, fmt.Errorf("convert: %w", err)
	}

	value, ok := mappedIssueTypes[jiraIssueType.Id]
	if !ok {
		return branch.NULL, errors.New("mapped issue type does not exist")
	}

	return value, nil
}
