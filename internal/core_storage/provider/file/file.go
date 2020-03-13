package file

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

var (
	bucketKV    = []byte("kv")
	bucketAlert = []byte("alert")
)

type Storage struct {
	name   string
	db     *bbolt.DB
	logger *zap.Logger
}

func New(config config.StorageCoreFile, logger *zap.Logger) (*Storage, error) {
	s := &Storage{
		name:   "file." + config.Name,
		logger: logger,
	}

	var err error

	options := &bbolt.Options{
		Timeout: config.Timeout,
	}

	s.db, err = bbolt.Open(config.Path, 0666, options)
	if err != nil {
		return nil, fmt.Errorf("error open db file, %w", err)
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("error init db, %w", err)
	}

	return s, nil
}

func (s *Storage) init() error {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketKV)
		if b == nil {
			_, err := tx.CreateBucket(bucketKV)
			if err != nil {
				return fmt.Errorf("error create bucket 'kv': %w", err)
			}
		}

		b = tx.Bucket(bucketAlert)
		if b == nil {
			_, err := tx.CreateBucket(bucketAlert)
			if err != nil {
				return fmt.Errorf("error create bucket 'alert': %w", err)
			}
		}

		return nil
	})

	return err
}

func (s *Storage) Name() string {
	return s.name
}

func (s *Storage) Stop() error {
	return s.db.Close()
}
