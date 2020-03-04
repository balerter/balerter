package config

type Storages struct {
	S3 []StorageS3 `json:"s3" yaml:"s3"`
}

func (cfg Storages) SetDefaults() {
	for _, c := range cfg.S3 {
		c.SetDefaults()
	}
}

func (cfg Storages) Validate() error {
	for _, c := range cfg.S3 {
		if err := c.Validate(); err != nil {
			return err
		}
	}

	return nil
}
