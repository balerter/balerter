package file

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"time"
)

var (
	bucketKV    = []byte("kv")
	bucketAlert = []byte("alert")

	defaultTimeout = time.Second
)

type storageKV struct {
	db     *bbolt.DB
	logger *zap.Logger
}

type storageAlert struct {
	db     *bbolt.DB
	logger *zap.Logger
}

type Storage struct {
	name   string
	db     *bbolt.DB
	logger *zap.Logger
	kv     *storageKV
	alert  *storageAlert
}

func (s *Storage) KV() coreStorage.KV {
	return s.kv
}

func (s *Storage) Alert() coreStorage.Alert {
	return s.alert
}

func New(config config.StorageCoreFile, logger *zap.Logger) (*Storage, error) {
	s := &Storage{
		name:   "file." + config.Name,
		logger: logger,
		kv: &storageKV{
			logger: logger,
		},
		alert: &storageAlert{
			logger: logger,
		},
	}

	var err error

	options := &bbolt.Options{
		Timeout: config.Timeout,
	}

	if options.Timeout == 0 {
		options.Timeout = defaultTimeout
	}

	s.db, err = bbolt.Open(config.Path, 0666, options)
	if err != nil {
		return nil, fmt.Errorf("error open db file, %w", err)
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("error init db, %w", err)
	}

	s.kv.db = s.db
	s.alert.db = s.db

	return s, nil
}

func (s *Storage) init() error {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketKV)
		if err != nil {
			return fmt.Errorf("error create bucket 'kv': %w", err)
		}

		_, err = tx.CreateBucketIfNotExists(bucketAlert)
		if err != nil {
			return fmt.Errorf("error create bucket 'alert': %w", err)
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
