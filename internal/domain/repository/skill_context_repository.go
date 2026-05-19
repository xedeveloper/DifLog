package repository

import "github.com/xedeveloper/DifLog/internal/domain/entity"

type SkillContextRepository interface {
	SaveSkillContext(ctx entity.SkillContext) error
	ListSkillContexts() ([]entity.SkillContext, error)
	LoadSkillContext(hash string) (entity.SkillContext, error)
}
