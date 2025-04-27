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

func trySaveConfig() error {
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
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

func GetString(token string) string {
    return viper.GetString(token)
}

func GetStringArray(token string) []string {
    return viper.GetStringSlice(token)
}

func GetSectionStringMap(section string) map[string][]string {
    return viper.GetStringMapStringSlice(section)
}

func GetAll() map[string]any {
    return viper.AllSettings()
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
