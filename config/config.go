package config

import (
    "errors"
    "fmt"
    "github.com/mitchellh/go-homedir"
    "github.com/spf13/viper"
    "twig/log"
)

const (
    Name = "twig"
    Type = "toml"
    Path = "/.config/twig"
)

type Token int

const (
    Project Token = iota
    ProjectHost
    ProjectAuth
    ProjectEmail
    ProjectToken

    Branch
    BranchDefault
    BranchOrigin
    BranchExclude

    Mapping
    MappingBuild
    MappingChore
    MappingCi
    MappingDocs
    MappingFeat
    MappingFix
    MappingPref
    MappingRefactor
    MappingRevert
    MappingStyle
    MappingTest
)

type Config struct {
    Project struct {
        Host  string `mapstructure:"host"`
        Auth  string `mapstructure:"auth"`
        Email string `mapstructure:"email"`
        Token string `mapstructure:"token"`
    } `mapstructure:"project"`
    Branch struct {
        Default string   `mapstructure:"default"`
        Origin  string   `mapstructure:"origin"`
        Exclude []string `mapstructure:"exclude"`
    } `mapstructure:"branch"`
    Mapping struct {
        Build    []int `mapstructure:"build"`
        Chore    []int `mapstructure:"chore"`
        Ci       []int `mapstructure:"ci"`
        Docs     []int `mapstructure:"docs"`
        Feat     []int `mapstructure:"feat"`
        Fix      []int `mapstructure:"fix"`
        Pref     []int `mapstructure:"pref"`
        Refactor []int `mapstructure:"refactor"`
        Revert   []int `mapstructure:"revert"`
        Style    []int `mapstructure:"style"`
        Test     []int `mapstructure:"test"`
    } `mapstructure:"mapping"`
}

func trySaveConfig() error {
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("failed to save config: %w", err)
    }
    return nil
}

func SetString(token string, value string) error {
    viper.Set(token, value)
    return trySaveConfig()
}

func SetStringArray(token string, value []string) error {
    viper.Set(token, value)
    return trySaveConfig()
}

func SetIntArray(token string, value []int) error {
    for _, v := range value {
        if v < 0 {
            return errors.New("values less than 0 are not permitted")
        }
    }

    viper.Set(token, value)
    return trySaveConfig()
}

func GetString(token Token) string {
    return viper.GetString(FromToken(token))
}

func GetStringArray(token Token) []string {
    return viper.GetStringSlice(FromToken(token))
}

func GetStringMap(token Token) map[string][]string {
    return viper.GetStringMapStringSlice(FromToken(token))
}

func GetAll() map[string]any {
    return viper.AllSettings()
}

func GetConfigSnapshot() (Config, error) {
    var cfg Config

    if err := viper.UnmarshalExact(&cfg); err != nil {
        return Config{}, fmt.Errorf("failed to get config snapshot: %w", err)
    }

    return cfg, nil
}

func InitConfig(file string) {
    if file != "" {
        log.Debug().Println("Using custom config")

        viper.SetConfigFile(file)
        viper.SetConfigType(Type)
    } else {
        log.Debug().Println("Using default config")

        home, err := homedir.Dir()
        if err != nil {
            log.Fatal().Println(err)
        }

        viper.AddConfigPath(home + Path)
        viper.SetConfigName(Name)
        viper.SetConfigType(Type)
    }

    if err := viper.ReadInConfig(); err != nil {
        log.Fatal().Println(err)
    }

    log.Debug().Println("Config loaded")
}

func FromToken(token Token) string {
    switch token {
    case Project:
        return "project"
    case ProjectHost:
        return "project.host"
    case ProjectAuth:
        return "project.auth"
    case ProjectEmail:
        return "project.email"
    case ProjectToken:
        return "project.token"
    case Branch:
        return "branch"
    case BranchDefault:
        return "branch.default"
    case BranchOrigin:
        return "branch.origin"
    case BranchExclude:
        return "branch.exclude"
    case Mapping:
        return "mapping"
    case MappingBuild:
        return "mapping.build"
    case MappingChore:
        return "mapping.chore"
    case MappingCi:
        return "mapping.ci"
    case MappingDocs:
        return "mapping.docs"
    case MappingFeat:
        return "mapping.feat"
    case MappingFix:
        return "mapping.fix"
    case MappingPref:
        return "mapping.pref"
    case MappingRefactor:
        return "mapping.refactor"
    case MappingRevert:
        return "mapping.revert"
    case MappingStyle:
        return "mapping.style"
    case MappingTest:
        return "mapping.test"
    default:
        return ""
    }
}
