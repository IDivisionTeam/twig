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

func castStringToInterface(src string) any {
    var tgt any
    tgt = src
    return tgt
}

func castSliceToInterface(src []string) []any {
    var tgt []any
    for k, v := range src {
        tgt[k] = v
    }
    return tgt
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
