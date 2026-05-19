package repository

import "github.com/xedeveloper/DifLog/internal/domain/entity"

type StagingRepository interface {
	LoadStagingIndex() (entity.StagingIndex, error)
	SaveStagingIndex(index entity.StagingIndex) error
	ClearStagingIndex() error
}
