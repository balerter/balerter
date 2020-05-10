package file

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

var (
	errBucketNotFound = errors.New("bucket not found")
)

func (s *storageKV) Put(key string, value string) error {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketKV)
		if b == nil {
			return errBucketNotFound
		}
		v := b.Get([]byte(key))
		if len(v) > 0 {
			return fmt.Errorf("key already exists")
		}

		return b.Put([]byte(key), []byte(value))
	})

	if err != nil {
		s.logger.Error("bbolt: error update item", zap.String("key", key), zap.String("value", value), zap.Error(err))
		return err
	}

	return nil
}

func (s *storageKV) All() (map[string]string, error) {
	result := make(map[string]string)

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketKV)
		if b == nil {
			return errBucketNotFound
		}

		err := b.ForEach(func(k, v []byte) error {
			result[string(k)] = string(v)
			return nil
		})
		if err != nil {
			return fmt.Errorf("error read storage")
		}

		return nil
	})

	if err != nil {
		s.logger.Error("bbolt: error get data", zap.Error(err))
		return nil, fmt.Errorf("error get data, %w", err)
	}

	return result, nil
}

func (s *storageKV) Get(key string) (string, error) {
	var v []byte

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketKV)
		if b == nil {
			return errBucketNotFound
		}
		v = b.Get([]byte(key))

		if len(v) == 0 {
			return fmt.Errorf("key not found")
		}

		return nil
	})

	if err != nil {
		s.logger.Error("bbolt: error get item", zap.String("key", key), zap.Error(err))
		return "", fmt.Errorf("error get item, %w", err)
	}

	return string(v), nil
}

func (s *storageKV) Upsert(key string, value string) error {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketKV)
		if b == nil {
			return errBucketNotFound
		}
		return b.Put([]byte(key), []byte(value))
	})

	if err != nil {
		s.logger.Error("bbolt: error update item", zap.String("key", key), zap.String("value", value), zap.Error(err))
		return fmt.Errorf("error update item, %w", err)
	}

	return nil
}

func (s *storageKV) Delete(key string) error {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketKV)
		if b == nil {
			return errBucketNotFound
		}
		return b.Delete([]byte(key))
	})

	if err != nil {
		s.logger.Error("bbolt: error delete item", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("error delete item, %w", err)
	}

	return nil
}
