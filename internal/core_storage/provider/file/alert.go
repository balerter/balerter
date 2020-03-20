package file

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

func (s *storageAlert) GetOrNew(name string) (*alert.Alert, error) {
	var v []byte

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketAlert)
		if b == nil {
			return errBucketNotFound
		}
		v = b.Get([]byte(name))

		return nil
	})

	if err != nil {
		s.logger.Error("bbolt: error get item", zap.ByteString("bucket", bucketAlert), zap.String("key", name), zap.Error(err))
		return nil, fmt.Errorf("error get item, %w", err)
	}

	a := alert.AcquireAlert()
	a.SetName(name)

	// if the buffer is empty, returns a new alert
	if len(v) == 0 {
		return a, nil
	}

	err = a.Unmarshal(v)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal alert, %w", err)
	}

	return a, nil

}

func (s *storageAlert) All() ([]*alert.Alert, error) {
	var vv [][]byte

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketAlert)
		if b == nil {
			return errBucketNotFound
		}
		err := b.ForEach(func(k, v []byte) error {
			vv = append(vv, v)
			return nil
		})

		return err
	})

	if err != nil {
		s.logger.Error("bbolt: error get items", zap.ByteString("bucket", bucketAlert), zap.Error(err))
		return nil, fmt.Errorf("error get item, %w", err)
	}

	res := make([]*alert.Alert, 0)

	for _, v := range vv {
		if len(v) == 0 {
			continue
		}

		a := alert.AcquireAlert()

		err = a.Unmarshal(v)
		if err != nil {
			return nil, fmt.Errorf("error unmarshal alert, %w", err)
		}

		res = append(res, a)
	}

	return res, nil
}

func (s *storageAlert) Release(a *alert.Alert) {
	alert.ReleaseAlert(a)
}
