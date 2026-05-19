package diffing_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xedeveloper/DifLog/internal/infrastructure/diffing"
)

func TestTextDifferService_DiffContents(t *testing.T) {
	svc := diffing.NewTextDifferService()

	diff := svc.DiffContents("hello world", "hello Go")
	assert.NotEmpty(t, diff)
}

func TestTextDifferService_UnifiedDiff(t *testing.T) {
	svc := diffing.NewTextDifferService()

	diff := svc.UnifiedDiff("old/file.md", "new/file.md", "line1\nline2\n", "line1\nline3\n")

	assert.Contains(t, diff, "--- old/file.md")
	assert.Contains(t, diff, "+++ new/file.md")
	assert.NotEmpty(t, diff, "diff should not be empty")
	_ = strings.Contains // used for import
}
