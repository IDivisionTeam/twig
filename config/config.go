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
    manager *viper.Viper
    name    string
    ext     string
    path    string
    homeDir string
}

var c *Config

func init() {
    c = New()
}

func New() *Config {
    c := new(Config)

    c.manager = viper.GetViper()
    c.name = "twig"
    c.ext = "toml"
    c.path = "/.config/twig"

    dir, err := homedir.Dir()
    if err != nil {
        log.Fatal().Println(err.Error())
    }

    c.homeDir = dir

    return c
}

func SetString(token Token, value string) error {
    return c.SetString(token, value)
}

func (c *Config) SetString(token Token, value string) error {
    key := FromToken(token)
    c.manager.Set(key, value)

    return c.overrideConfig()
}

func SetStringArray(token Token, value []string) error {
    return c.SetStringArray(token, value)
}

func (c *Config) SetStringArray(token Token, value []string) error {
    key := FromToken(token)
    c.manager.Set(key, value)

    return c.overrideConfig()
}

func GetString(token Token) string {
    return c.GetString(token)
}

func (c *Config) GetString(token Token) string {
    key := FromToken(token)
    return c.manager.GetString(key)
}

func GetStringArray(token Token) []string {
    return c.GetStringArray(token)
}

func (c *Config) GetStringArray(token Token) []string {
    key := FromToken(token)
    return c.manager.GetStringSlice(key)
}

func GetStringMap(token Token) map[string][]string {
    return c.GetStringMap(token)
}

func (c *Config) GetStringMap(token Token) map[string][]string {
    key := FromToken(token)
    return c.manager.GetStringMapStringSlice(key)
}

func GetAllSnapshot() (*Settings, error) {
    return c.GetAllSnapshot()
}

func (c *Config) GetAllSnapshot() (*Settings, error) {
    var cfg *Settings

    if err := c.manager.UnmarshalExact(&cfg); err != nil {
        return nil, fmt.Errorf("failed to get config snapshot: %w", err)
    }

    return cfg, nil
}

func InitConfig(file string) {
    if file != "" {
        log.Debug().Println("Using custom config")

        c.manager.SetConfigFile(file)
        c.manager.SetConfigType(c.ext)
    } else {
        log.Debug().Println("Using default config")

        c.manager.AddConfigPath(c.homeDir + c.path)
        c.manager.SetConfigName(c.name)
        c.manager.SetConfigType(c.ext)
    }

    if err := c.manager.ReadInConfig(); err != nil {
        log.Fatal().Println(err)
    }

    log.Debug().Println("Config loaded")
}

func IsConfigExist() error {
    fileName := fmt.Sprintf("%s.%s", c.name, c.ext)
    path := filepath.Join(c.homeDir, c.path, fileName)

    _, err := os.Stat(path)

    if err == nil || os.IsExist(err) {
        return errors.New("config exist, abort")
    }
    return nil
}

func CreateConfigDir() error {
    path := filepath.Join(c.homeDir, c.path)
    err := os.MkdirAll(path, os.ModeDir)

    if err == nil || os.IsExist(err) {
        return nil
    } else {
        return err
    }
}

func CreatConfigFile() error {
    fileName := fmt.Sprintf("%s.%s", c.name, c.ext)
    path := filepath.Join(c.homeDir, c.path, fileName)

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

func GetDefaultConfigPath() string {
    return fmt.Sprintf(
        "%s/%s.%s",
        c.path,
        c.name,
        c.ext,
    )
}

func (c *Config) overrideConfig() error {
    if err := c.manager.WriteConfig(); err != nil {
        return fmt.Errorf("failed to save config: %w", err)
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
