package repository

import "github.com/xedeveloper/DifLog/internal/domain/entity"

type ReleaseRepository interface{
	GetLatestRelease() (entity.Release, error)
}
