package bindata

import (
	"aletheia.icu/broccoli/fs"
)

//go:generate broccoli -src=../modules -o modules

func Broccoli() *fs.Broccoli {
	return br
}
