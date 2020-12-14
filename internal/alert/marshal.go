package alert

import (
	"encoding/binary"
	"errors"
)

const (
	timeFieldBinarySize = 15
)

var (
	ErrDecodeAlertName = errors.New("error decode alert name")
	ErrSourceTooSmall  = errors.New("source too small")
	ErrDecodeLevel     = errors.New("error decode level")
	ErrDecodeCount     = errors.New("error decode count")
	ErrSourceTooLong   = errors.New("source too long")
)

// Marshal an Alert to the byte slice
func (a *Alert) Marshal() ([]byte, error) {
	res := make([]byte, 0)

	buf := make([]byte, 64)
	var n int

	// Name
	n = binary.PutUvarint(buf, uint64(len(a.name)))
	res = append(res, buf[:n]...)
	res = append(res, a.name...)

	// Level
	n = binary.PutUvarint(buf, uint64(a.level))
	res = append(res, buf[:n]...)

	// LastChange
	t, err := a.lastChange.MarshalBinary()
	if err != nil {
		return nil, err
	}
	res = append(res, t...)

	// Start
	t, err = a.start.MarshalBinary()
	if err != nil {
		return nil, err
	}
	res = append(res, t...)

	// Count
	n = binary.PutUvarint(buf, uint64(a.count))
	res = append(res, buf[:n]...)

	return res, nil
}

// Unmarshal the byte slice to an Alert
func (a *Alert) Unmarshal(src []byte) error {
	// Name length
	l, n := binary.Uvarint(src)
	if n <= 0 {
		return ErrDecodeAlertName
	}
	src = src[n:]

	// Name
	if len(src) < int(l) {
		return ErrSourceTooSmall
	}
	a.name = string(src[:l])
	src = src[l:]

	// Level
	if len(src) == 0 {
		return ErrSourceTooSmall
	}
	l, n = binary.Uvarint(src)
	if n <= 0 {
		return ErrDecodeLevel
	}
	a.level = Level(l)
	src = src[n:]

	// LastChange
	if len(src) < timeFieldBinarySize {
		return ErrSourceTooSmall
	}
	err := a.lastChange.UnmarshalBinary(src[:timeFieldBinarySize])
	if err != nil {
		return err
	}
	src = src[timeFieldBinarySize:]

	// Start
	if len(src) < timeFieldBinarySize {
		return ErrSourceTooSmall
	}
	err = a.start.UnmarshalBinary(src[:timeFieldBinarySize])
	if err != nil {
		return err
	}
	src = src[timeFieldBinarySize:]

	// Count
	if len(src) == 0 {
		return ErrSourceTooSmall
	}
	l, n = binary.Uvarint(src)
	if n <= 0 {
		return ErrDecodeCount
	}
	a.count = int(l)
	src = src[n:]

	if len(src) > 0 {
		return ErrSourceTooLong
	}

	return nil
}
