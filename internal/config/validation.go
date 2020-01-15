package config

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func (cfg *Config) Validate() error {
	for _, c := range cfg.Scripts.Sources.Folder {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	for _, c := range cfg.DataSources.Clickhouse {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	for _, c := range cfg.DataSources.Prometheus {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	for _, c := range cfg.Channels.Slack {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *DataSourcePrometheus) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if strings.TrimSpace(cfg.URL) == "" {
		return fmt.Errorf("url must be not empty")
	}

	return nil
}

func (cfg *DataSourceClickhouse) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if cfg.Host == "" {
		return fmt.Errorf("host must be defined")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("port must be defined")
	}
	if cfg.Username == "" {
		return fmt.Errorf("username must be defined")
	}

	return nil
}

func (cfg *ScriptSourceFolder) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	_, err := ioutil.ReadDir(cfg.Path)
	if err != nil {
		return fmt.Errorf("error read folder '%s', %w", cfg.Path, err)
	}

	return nil
}

func (cfg *ChannelSlack) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}

	if strings.TrimSpace(cfg.URL) == "" {
		return fmt.Errorf("url must be not empty")
	}

	return nil
}
