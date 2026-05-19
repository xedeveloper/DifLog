package service

type DifferService interface {
	DiffContents(oldContent, newContent string) string
	UnifiedDiff(oldLabel, newLabel, oldContent, newContent string) string
}
