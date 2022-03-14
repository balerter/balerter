package prometheus_models

const (
	ValString  = "string"
	ValScalar  = "scalar"
	ValVector  = "vector"
	ValMatrix  = "matrix"
	ValStreams = "streams" // loki
)

type ModelValue interface {
	Type() string
	String() string
}
