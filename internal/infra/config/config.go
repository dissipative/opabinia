package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	LogLevel    string `yaml:"log_level"`
	CacheSize   int    `yaml:"cache_size"`
	CacheExpire int    `yaml:"cache_expire"`
	Env         string `yaml:"env"`
}

type Template struct {
	File        string `yaml:"file"`
	Description string `yaml:"description"`
	TitlePrefix string `yaml:"title_prefix"`

	*Paths
}

type Paths struct {
	CustomCSS   string `yaml:"custom_css"`
	CustomJS    string `yaml:"custom_js"`
	FaviconDir  string `yaml:"favicon_dir"`
	Webmanifest string `yaml:"webmanifest"`
}

type Markdown struct {
	CodeDefaultLang string `yaml:"code_default_lang"`
	SyntaxTheme     string `yaml:"syntax_theme"`
}

type Config struct {
	*Server
	*Template
	*Markdown
}

func NewConfig() *Config {
	return &Config{
		// defaults filled
		Server: &Server{
			Host:        "localhost",
			Port:        4321,
			LogLevel:    "debug",
			CacheSize:   0,
			CacheExpire: 3000,
			Env:         "development",
		},
		Template: &Template{
			File:  "templates/default.tmpl",
			Paths: &Paths{},
		},
		Markdown: &Markdown{
			CodeDefaultLang: "bash",
			SyntaxTheme:     "nord",
		},
	}
}

func Load() (*Config, error) {
	c := NewConfig()

	file, err := os.ReadFile("opabinia.yml")
	if err != nil {
		return c, nil
	}

	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) LogLevel() []byte {
	return []byte(c.Server.LogLevel)
}

func (c *Config) TmplConf() *Template {
	return c.Template
}

func (c *Config) SyntaxTheme() string {
	return c.Markdown.SyntaxTheme
}

func (c *Config) CodeDefaultLang() string {
	return c.Markdown.CodeDefaultLang
}

func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) IsDevEnv() bool {
	return c.Server.Env == "development" || c.Server.Env == "dev"
}
