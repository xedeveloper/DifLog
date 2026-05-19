package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

type SHA256HasherService struct{}

func NewSHA256HasherService() *SHA256HasherService {
	return &SHA256HasherService{}
}

func (s *SHA256HasherService) HashContent(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func (s *SHA256HasherService) HashMultiple(contents [][]byte) string {
	hashes := make([]string, len(contents))
	for i, c := range contents {
		hashes[i] = s.HashContent(c)
	}
	sort.Strings(hashes)
	combined := ""
	for _, h := range hashes {
		combined += h
	}
	return s.HashContent([]byte(combined))
}
