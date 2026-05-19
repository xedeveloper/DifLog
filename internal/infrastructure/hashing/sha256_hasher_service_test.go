package hashing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xedeveloper/DifLog/internal/infrastructure/hashing"
)

func TestSHA256HasherService_HashContent(t *testing.T) {
	svc := hashing.NewSHA256HasherService()

	hash1 := svc.HashContent([]byte("hello world"))
	hash2 := svc.HashContent([]byte("hello world"))
	hash3 := svc.HashContent([]byte("different content"))

	assert.Equal(t, hash1, hash2, "same content should produce same hash")
	assert.NotEqual(t, hash1, hash3, "different content should produce different hash")
	assert.Len(t, hash1, 64, "SHA256 hex should be 64 chars")
}

func TestSHA256HasherService_HashMultiple(t *testing.T) {
	svc := hashing.NewSHA256HasherService()

	contents := [][]byte{[]byte("file1"), []byte("file2")}
	hash := svc.HashMultiple(contents)

	assert.Len(t, hash, 64)

	// Order-independent: same files in different order should produce same hash
	reversed := [][]byte{[]byte("file2"), []byte("file1")}
	hashReversed := svc.HashMultiple(reversed)
	assert.Equal(t, hash, hashReversed, "hash should be order-independent")
}
