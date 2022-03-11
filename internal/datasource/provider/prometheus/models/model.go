package models

const (
	ValString = "string"
	ValScalar = "scalar"
	ValVector = "vector"
	ValMatrix = "matrix"
)

type ModelValue interface {
	Type() string
	String() string
}
