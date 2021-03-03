package loki

import (
	"fmt"
	"github.com/balerter/balerter/internal/config/common"
	"strings"
	"time"
)

type Loki struct {
	Name      string            `json:"name" yaml:"name" hcl:"name,label"`
	URL       string            `json:"url" yaml:"url" hcl:"url"`
	BasicAuth *common.BasicAuth `json:"basicAuth" yaml:"basicAuth" hcl:"basicAuth,block"`
	Timeout   time.Duration     `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
}

func (cfg Loki) Validate() error {
	if strings.TrimSpace(cfg.Name) == "" {
		return fmt.Errorf("name must be not empty")
	}
	if strings.TrimSpace(cfg.URL) == "" {
		return fmt.Errorf("url must be not empty")
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	return nil
}
