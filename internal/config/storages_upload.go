package config

type StoragesUpload struct {
	S3 []StorageUploadS3 `json:"s3" yaml:"s3"`
}

func (cfg StoragesUpload) SetDefaults() {
	for _, c := range cfg.S3 {
		c.SetDefaults()
	}
}

func (cfg StoragesUpload) Validate() error {
	for _, c := range cfg.S3 {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}
