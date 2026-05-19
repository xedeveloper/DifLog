package service

type HasherService interface {
	HashContent(content []byte) string
	HashMultiple(contents [][]byte) string
}
