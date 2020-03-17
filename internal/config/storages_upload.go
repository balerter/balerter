package config

import "fmt"

type StoragesUpload struct {
	S3 []StorageUploadS3 `json:"s3" yaml:"s3"`
}

func (cfg StoragesUpload) Validate() error {
	var names []string

	for _, c := range cfg.S3 {
		names = append(names, c.Name)
		if err := c.Validate(); err != nil {
			return err
		}
	}
	if name := checkUnique(names); name != "" {
		return fmt.Errorf("found duplicated name for upload storages 's3': %s", name)
	}

	return nil
}
