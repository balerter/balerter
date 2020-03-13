package core_storage

type CoreStorage interface {
	Name() string
	Put(string, string) error
	Get(string) (string, error)
	Upsert(string, string) error
	Delete(string) error
}
