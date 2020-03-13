package alert

import (
	"encoding/binary"
	"fmt"
)

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

func (a *Alert) Unmarshal(src []byte) error {

	// Name length
	l, n := binary.Uvarint(src)
	if n <= 0 {
		return fmt.Errorf("error decode alert name")
	}
	src = src[n:]

	// Name
	if len(src) < int(l) {
		return fmt.Errorf("source too small")
	}
	a.name = string(src[:l])
	src = src[l:]

	// Level
	if len(src) == 0 {
		return fmt.Errorf("source too small")
	}
	l, n = binary.Uvarint(src)
	if n <= 0 {
		return fmt.Errorf("error decode level")
	}
	a.level = Level(l)
	src = src[n:]

	// LastChange
	if len(src) < 15 {
		return fmt.Errorf("source too small")
	}
	err := a.lastChange.UnmarshalBinary(src[:15])
	if err != nil {
		return err
	}
	src = src[15:]

	// Start
	if len(src) < 15 {
		return fmt.Errorf("source too small")
	}
	err = a.start.UnmarshalBinary(src[:15])
	if err != nil {
		return err
	}
	src = src[15:]

	// Count
	if len(src) == 0 {
		return fmt.Errorf("source too small")
	}
	l, n = binary.Uvarint(src)
	if n <= 0 {
		return fmt.Errorf("error decode last change")
	}
	a.count = int(l)
	src = src[n:]

	if len(src) > 0 {
		return fmt.Errorf("source too long")
	}

	return nil
}
