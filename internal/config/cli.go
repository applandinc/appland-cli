package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type Config struct {
	CurrentContext string              `yaml:"current_context"`
	Contexts       map[string]*Context `yaml:"contexts"`
}

type Context struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
}

const (
	applandFilename    string = ".appland"
	defaultContextName string = "default"
)

var (
	config         *Config
	configPath     string
	currentContext *Context
	defaultContext = Context{
		URL:    "https://app.land",
		APIKey: "",
	}
)

func makeDefault() *Config {
	return &Config{
		CurrentContext: defaultContextName,
		Contexts: map[string]*Context{
			defaultContextName: &defaultContext,
		},
	}
}

func loadCLIConfig(path string) bool {
	info, err := getFS().Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "warn: %s\n", err)
		return false
	}

	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "warn: %s is a directory\n", path)
		return false
	}

	data, err := afero.ReadFile(getFS(), path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warn: %s\v", err)
		return false
	}

	c := &Config{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warn: %s\n", err)
		return false
	}

	config = c
	configPath = path

	if len(config.Contexts) == 0 {
		config = makeDefault()
	}

	return true
}

func LoadCLIConfig() {
	envPath := os.Getenv("APPLAND_CONFIG")
	if envPath != "" && loadCLIConfig(envPath) {
		return
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warn: %s\n", err)
	} else {
		if loadCLIConfig(path.Join(currentDir, applandFilename)) {
			return
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warn: %s\n", err)
	} else {
		if loadCLIConfig(path.Join(homeDir, applandFilename)) {
			return
		}
	}

	config = makeDefault()
	configPath = path.Join(homeDir, applandFilename)
}

func WriteCLIConfig() error {
	if configPath == "" {
		return fmt.Errorf("no config path is set")
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return afero.WriteFile(getFS(), configPath, data, 0600)
}

func GetAPIKey() string {
	if currentContext == nil {
		return ""
	}

	return currentContext.APIKey
}

func SetCurrentContext(name string) error {
	if GetContext(name) == nil {
		return fmt.Errorf("no context named '%s' exists", name)
	}

	config.CurrentContext = name
	return nil
}

func SetContextUrl(url string) {
	currentContext.URL = url
}

func SetContextAPIKey(apiKey string) {
	currentContext.APIKey = apiKey
}

func GetContext(name string) *Context {
	return config.Contexts[name]
}

func GetCurrentContext() *Context {
	return GetContext(config.CurrentContext)
}

func GetCurrentContextName() string {
	return config.CurrentContext
}

func RenameContext(old string, new string) {
	if old == new {
		return
	}

	config.Contexts[new] = config.Contexts[old]
	delete(config.Contexts, old)

	if config.CurrentContext == old {
		config.CurrentContext = new
	}
}

func MakeContext(name string, url string) error {
	if GetContext(name) != nil {
		return fmt.Errorf("a context named '%s' already exists", name)
	}

	config.Contexts[name] = &Context{
		URL: url,
	}

	return nil
}

func GetCLIConfig() *Config {
	return config
}
