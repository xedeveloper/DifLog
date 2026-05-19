package repository

type ContextRepository interface {
	SaveObject(hash string, content []byte) error
	LoadObject(hash string) ([]byte, error)
	ObjectExists(hash string) bool
}
