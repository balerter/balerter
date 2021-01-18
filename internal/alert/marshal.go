package alert

import (
	"encoding/binary"
	"errors"
)

const (
	timeFieldBinarySize = 15
)

var (
	ErrDecodeAlertName = errors.New("error decode alert Name")
	ErrSourceTooSmall  = errors.New("source too small")
	ErrDecodeLevel     = errors.New("error decode Level")
	ErrDecodeCount     = errors.New("error decode Count")
	ErrSourceTooLong   = errors.New("source too long")
)

// Marshal an Alert to the byte slice
func (a *Alert) Marshal() ([]byte, error) {
	res := make([]byte, 0)

	buf := make([]byte, 64)
	var n int

	// Name
	n = binary.PutUvarint(buf, uint64(len(a.Name)))
	res = append(res, buf[:n]...)
	res = append(res, a.Name...)

	// Level
	n = binary.PutUvarint(buf, uint64(a.Level))
	res = append(res, buf[:n]...)

	// LastChange
	t, err := a.LastChange.MarshalBinary()
	if err != nil {
		return nil, err
	}
	res = append(res, t...)

	// Start
	t, err = a.Start.MarshalBinary()
	if err != nil {
		return nil, err
	}
	res = append(res, t...)

	// Count
	n = binary.PutUvarint(buf, uint64(a.Count))
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
	a.Name = string(src[:l])
	src = src[l:]

	// Level
	if len(src) == 0 {
		return ErrSourceTooSmall
	}
	l, n = binary.Uvarint(src)
	if n <= 0 {
		return ErrDecodeLevel
	}
	a.Level = Level(l)
	src = src[n:]

	// LastChange
	if len(src) < timeFieldBinarySize {
		return ErrSourceTooSmall
	}
	err := a.LastChange.UnmarshalBinary(src[:timeFieldBinarySize])
	if err != nil {
		return err
	}
	src = src[timeFieldBinarySize:]

	// Start
	if len(src) < timeFieldBinarySize {
		return ErrSourceTooSmall
	}
	err = a.Start.UnmarshalBinary(src[:timeFieldBinarySize])
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
	a.Count = int(l)
	src = src[n:]

	if len(src) > 0 {
		return ErrSourceTooLong
	}

	return nil
}
