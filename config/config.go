package config

import (
    _ "embed"
    "errors"
    "fmt"
    "github.com/mitchellh/go-homedir"
    "github.com/spf13/viper"
    "os"
    "path/filepath"
    "twig/log"
)

//go:embed twig.toml
var embededCfg []byte

const (
    Name = "twig"
    Type = "toml"
    Path = "/.config/twig"
)

type Token int

const (
    Unspecified Token = iota - 1
    Project
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
    Project ProjectCfg `mapstructure:"project"`
    Branch  BranchCfg  `mapstructure:"branch"`
    Mapping MappingCfg `mapstructure:"mapping"`
}

type ProjectCfg struct {
    Host  string `mapstructure:"host"`
    Auth  string `mapstructure:"auth"`
    Email string `mapstructure:"email"`
    Token string `mapstructure:"token"`
}

type BranchCfg struct {
    Default string   `mapstructure:"default"`
    Origin  string   `mapstructure:"origin"`
    Exclude []string `mapstructure:"exclude"`
}

type MappingCfg struct {
    Build    []string `mapstructure:"build"`
    Chore    []string `mapstructure:"chore"`
    Ci       []string `mapstructure:"ci"`
    Docs     []string `mapstructure:"docs"`
    Feat     []string `mapstructure:"feat"`
    Fix      []string `mapstructure:"fix"`
    Pref     []string `mapstructure:"pref"`
    Refactor []string `mapstructure:"refactor"`
    Revert   []string `mapstructure:"revert"`
    Style    []string `mapstructure:"style"`
    Test     []string `mapstructure:"test"`
}

func overrideConfig() error {
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("failed to save config: %w", err)
    }
    return nil
}

func SetString(token Token, value string) error {
    viper.Set(FromToken(token), value)
    return overrideConfig()
}

func SetStringArray(token Token, value []string) error {
    viper.Set(FromToken(token), value)
    return overrideConfig()
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

func IsConfigExist() error {
    home, err := homedir.Dir()
    if err != nil {
        log.Fatal().Println(err)
    }

    fileName := fmt.Sprintf("%s.%s", Name, Type)
    path := filepath.Join(home, Path, fileName)

    _, err = os.Stat(path)

    if err == nil || os.IsExist(err) {
        return errors.New("config exist, abort")
    }
    return nil
}

func CreateConfigDir() error {
    home, err := homedir.Dir()
    if err != nil {
        log.Fatal().Println(err)
    }

    path := filepath.Join(home, Path)
    err = os.MkdirAll(path, os.ModeDir)

    if err == nil || os.IsExist(err) {
        return nil
    } else {
        return err
    }
}

func CreatConfigFile() error {
    home, err := homedir.Dir()
    if err != nil {
        log.Fatal().Println(err)
    }

    fileName := fmt.Sprintf("%s.%s", Name, Type)
    path := filepath.Join(home, Path, fileName)

    file, err := os.Create(path)
    if err != nil {
        return err
    }

    _, err = file.Write(embededCfg)
    if err != nil {
        return err
    }

    if err = file.Close(); err != nil {
        return err
    }

    return nil
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

func FromInput(token string) (Token, error) {
    switch token {
    case "project":
        return Project, nil
    case "project.host":
        return ProjectHost, nil
    case "project.auth":
        return ProjectAuth, nil
    case "project.email":
        return ProjectEmail, nil
    case "project.token":
        return ProjectToken, nil
    case "branch":
        return Branch, nil
    case "branch.default":
        return BranchDefault, nil
    case "branch.origin":
        return BranchOrigin, nil
    case "branch.exclude":
        return BranchExclude, nil
    case "mapping":
        return Mapping, nil
    case "mapping.build":
        return MappingBuild, nil
    case "mapping.chore":
        return MappingChore, nil
    case "mapping.ci":
        return MappingCi, nil
    case "mapping.docs":
        return MappingDocs, nil
    case "mapping.feat":
        return MappingFeat, nil
    case "mapping.fix":
        return MappingFix, nil
    case "mapping.pref":
        return MappingPref, nil
    case "mapping.refactor":
        return MappingRefactor, nil
    case "mapping.revert":
        return MappingRevert, nil
    case "mapping.style":
        return MappingStyle, nil
    case "mapping.test":
        return MappingTest, nil
    default:
        return Unspecified, errors.New("unexpected token from input")
    }
}
