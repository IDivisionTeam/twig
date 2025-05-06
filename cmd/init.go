package cmd

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"twig/config"
	"twig/log"
	"twig/network"
)

var (
	InitCmdName = "init"
	initCmd = &cobra.Command{
		Use:   InitCmdName,
		Short: "Create config",
		Args:  cobra.NoArgs,
		Run:   runInit,
	}
)

func runInit(cmd *cobra.Command, args []string) {
	input := prompt.NewStandardInputParser()
	c := color.New(color.FgHiGreen)
	cmdName := cmd.Name()

	log.Debug().Println(fmt.Sprintf("%s: executing command", cmdName))
	log.Warn().Println("Press \"ENTER\" to keep the default value")

	if err := config.CreateConfigDir(); err != nil {
		logCmdFatal(err)
	}

	if err := config.IsConfigExist(); err == nil {
		if err = config.CreatConfigFile(); err != nil {
			logCmdFatal(err)
		}
	}

	config.InitConfig("")

	log.Info().Println("\nProject Group")
	if err := setHostFromInput(c, input); err != nil {
		logCmdFatal(err)
	}

	if err := setEmailFromInput(c, input); err != nil {
		logCmdFatal(err)
	}

	if err := setTokenFromInput(c, input); err != nil {
		logCmdFatal(err)
	}

	if err := setAuthFromInput(c, input); err != nil {
		logCmdFatal(err)
	}

	log.Info().Println("\nBranch Group")
	if err := setBranchDefaultFromInput(c, input); err != nil {
		logCmdFatal(err)
	}

	if err := setBranchOriginFromInput(c, input); err != nil {
		logCmdFatal(err)
	}

	if err := setBranchExcludesFromInput(c, input); err != nil {
		logCmdFatal(err)
	}

	log.Info().Println("\nMapping Group")
	if err := setMappingFromInput(c, input); err != nil {
		logCmdFatal(err)
	}

	log.Info().Println("\nSetup complete. You're ready to go")
}

func setHostFromInput(c *color.Color, in *prompt.PosixParser) error {
	fmt.Print(c.Sprint("What is your JIRA host? (e.g. example.atlassian.net): "))

	str, err := in.Read()
	if err != nil {
		return err
	}

	value := strings.TrimSpace(string(str))
	if value == "" {
		return nil // skip, using default
	}
	if err = config.SetString(config.ProjectHost, value); err != nil {
		return err
	}
	log.Debug().Println(fmt.Sprintf("Input host: %q", value))

	return nil
}

func setEmailFromInput(c *color.Color, in *prompt.PosixParser) error {
	fmt.Print(c.Sprint("What is your email in Jira? (e.g. example@exp.com): "))

	str, err := in.Read()
	if err != nil {
		return err
	}

	value := strings.TrimSpace(string(str))
	if value == "" {
		return nil // skip, using default
	}
	if err = config.SetString(config.ProjectEmail, value); err != nil {
		return err
	}
	log.Debug().Println(fmt.Sprintf("Input email: %q", value))

	return nil
}

func setTokenFromInput(c *color.Color, in *prompt.PosixParser) error {
	fmt.Print(c.Sprint("What is your Jira PAT? (e.g. dXNlckBleGFtcGx): "))

	str, err := in.Read()
	if err != nil {
		return err
	}

	value := strings.TrimSpace(string(str))
	if value == "" {
		return nil // skip, using default
	}
	if err = config.SetString(config.ProjectToken, value); err != nil {
		return err
	}
	log.Debug().Println(fmt.Sprintf("Input token: %q", value))

	return nil
}

func setAuthFromInput(c *color.Color, in *prompt.PosixParser) error {
	for {
		fmt.Print(c.Sprint("What is your Jira Auth type? (basic/bearer): "))

		str, err := in.Read()
		if err != nil {
			return err
		}

		value := strings.ToLower(string(str))
		value = strings.TrimSpace(value)

		switch value {
		case network.BasicType, network.BearerType:
			if err = config.SetString(config.ProjectAuth, value); err != nil {
				return err
			}
			log.Debug().Println(fmt.Sprintf("Input auth: %q", value))
			return nil
		case "":
			return nil // skip, using default
		default:
			msg := fmt.Sprintf("Invalid input. Please enter %q for Basic or %q for Bearer", network.BasicType, network.BearerType)
			log.Error().Println(msg)
			continue
		}
	}
}

func setBranchDefaultFromInput(c *color.Color, in *prompt.PosixParser) error {
	fmt.Print(c.Sprint("What branch do you use as default? (e.g. development): "))

	str, err := in.Read()
	if err != nil {
		return err
	}
	value := strings.TrimSpace(string(str))
	if value == "" {
		return nil // skip, using default
	}
	if err = config.SetString(config.BranchDefault, value); err != nil {
		return err
	}
	log.Debug().Println(fmt.Sprintf("Input branch: %q", value))

	return nil
}

func setBranchOriginFromInput(c *color.Color, in *prompt.PosixParser) error {
	fmt.Print(c.Sprint("What remote do you use as default? (e.g. origin): "))

	str, err := in.Read()
	if err != nil {
		return err
	}
	value := strings.TrimSpace(string(str))
	if value == "" {
		return nil // skip, using default
	}
	if err = config.SetString(config.BranchOrigin, value); err != nil {
		return err
	}
	log.Debug().Println(fmt.Sprintf("Input origin: %q", value))

	return nil
}

func setBranchExcludesFromInput(c *color.Color, in *prompt.PosixParser) error {
	fmt.Print(c.Sprint("Exclude any words from the branch name? (e.g. be,mobile,web): "))

	str, err := in.Read()
	if err != nil {
		return err
	}
	value := strings.TrimSpace(string(str))
	if value == "" {
		return nil // skip, using default
	}

	valueArr := strings.Split(value, ",")
	for i, _ := range valueArr {
		// and trim space just in case
		valueArr[i] = strings.TrimSpace(valueArr[i])
	}
	if err = config.SetStringArray(config.BranchExclude, valueArr); err != nil {
		return err
	}
	log.Debug().Println(fmt.Sprintf("Input exclude: %q", value))

	return nil
}

func setMappingFromInput(c *color.Color, in *prompt.PosixParser) error {
	log.Info().Println("Please input a valid id/ids for the following options:\nbuild, chore, ci, docs, feat, fix, pref, refactor, revert, style, test")

	start := config.Mapping + 1
	end := config.Mapping + 11

	for {
		if start > end {
			return nil
		}

		option := strings.SplitAfter(config.FromToken(start), ".")[1]
		fmt.Print(c.Sprintf("What should %q be mapped to? (e.g. 101 or 101,102): ", option))

		str, err := in.Read()
		if err != nil {
			return err
		}

		value := strings.TrimSpace(string(str))
		if value == "" {
			start++
			continue // skip, using default
		}

		valueArr := strings.Split(value, ",")
		for i, _ := range valueArr {
			// and trim space just in case
			valueArr[i] = strings.TrimSpace(valueArr[i])
		}

		hasIncorrectVal := false
		for _, v := range valueArr {
			// must check if input correct
			if _, err = strconv.Atoi(v); err != nil {
				hasIncorrectVal = true
				break
			}
		}

		if hasIncorrectVal {
			log.Error().Println("Invalid input. Please enter integer value or integer values split by comma")
			continue
		}

		if err = config.SetStringArray(start, valueArr); err != nil {
			log.Error().Println(err.Error())
		}
		start++
	}
}
