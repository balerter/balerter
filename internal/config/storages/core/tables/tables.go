package tables

import "fmt"

// AlertFields describe alerts table fields
type AlertFields struct {
	Name      string `json:"name" yaml:"name" hcl:"name"`
	Level     string `json:"level" yaml:"level" hcl:"level"`
	Count     string `json:"count" yaml:"count" hcl:"count"`
	UpdatedAt string `json:"updatedAt" yaml:"updatedAt" hcl:"updatedAt"`
	CreatedAt string `json:"createdAt" yaml:"createdAt" hcl:"createdAt"`
}

// KVFields describe KV table fields
type KVFields struct {
	Key   string `json:"key" yaml:"key" hcl:"key"`
	Value string `json:"value" yaml:"value" hcl:"value"`
}

// TableKV is config for core storage kv table
type TableKV struct {
	Table       string   `json:"table" yaml:"table" hcl:"table"`
	Fields      KVFields `json:"fields" yaml:"fields" hcl:"fields,block"`
	CreateTable bool     `json:"create" yaml:"create" hcl:"create,optional"`
}

// TableAlerts is config for core storage alerts table
type TableAlerts struct {
	Table       string      `json:"table" yaml:"table" hcl:"table"`
	Fields      AlertFields `json:"fields" yaml:"fields" hcl:"fields,block"`
	CreateTable bool        `json:"create" yaml:"create" hcl:"create,optional"`
}

// Validate config
func (t TableAlerts) Validate() error {
	if t.Table == "" {
		return fmt.Errorf("table must be not empty")
	}
	if err := t.Fields.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate config
func (t TableKV) Validate() error {
	if t.Table == "" {
		return fmt.Errorf("table must be not empty")
	}
	if err := t.Fields.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate config
func (t AlertFields) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("field name must be not empty")
	}
	if t.Level == "" {
		return fmt.Errorf("field level must be not empty")
	}
	if t.Count == "" {
		return fmt.Errorf("field count must be not empty")
	}
	if t.CreatedAt == "" {
		return fmt.Errorf("field createdAt must be not empty")
	}
	if t.UpdatedAt == "" {
		return fmt.Errorf("field updatedAt must be not empty")
	}
	return nil
}

// Validate config
func (t KVFields) Validate() error {
	if t.Key == "" {
		return fmt.Errorf("field key must be not empty")
	}
	if t.Value == "" {
		return fmt.Errorf("field key must be not empty")
	}
	return nil
}
