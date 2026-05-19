package diffing

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type TextDifferService struct {
	dmp *diffmatchpatch.DiffMatchPatch
}

func NewTextDifferService() *TextDifferService {
	return &TextDifferService{dmp: diffmatchpatch.New()}
}

func (t *TextDifferService) DiffContents(oldContent, newContent string) string {
	diffs := t.dmp.DiffMain(oldContent, newContent, false)
	t.dmp.DiffCleanupSemantic(diffs)
	return t.dmp.DiffPrettyText(diffs)
}

func (t *TextDifferService) UnifiedDiff(oldLabel, newLabel, oldContent, newContent string) string {
	diffs := t.dmp.DiffMain(oldContent, newContent, false)
	t.dmp.DiffCleanupSemantic(diffs)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("--- %s\n", oldLabel))
	sb.WriteString(fmt.Sprintf("+++ %s\n", newLabel))

	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		for _, line := range lines {
			switch diff.Type {
			case diffmatchpatch.DiffInsert:
				sb.WriteString(fmt.Sprintf("+ %s\n", line))
			case diffmatchpatch.DiffDelete:
				sb.WriteString(fmt.Sprintf("- %s\n", line))
			case diffmatchpatch.DiffEqual:
				sb.WriteString(fmt.Sprintf("  %s\n", line))
			}
		}
	}
	return sb.String()
}
