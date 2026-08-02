package remote

import (
	"context"

	"github.com/google/go-github/v66/github"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

const (
	releaseOwner = "xedeveloper"
	releaseRepo  = "DifLog"
)

type GithubReleaseRepository struct {
	client *github.Client
}

func NewGithubReleaseRepository() *GithubReleaseRepository {
	return &GithubReleaseRepository{client: github.NewClient(nil)}
}

func (r *GithubReleaseRepository) GetLatestRelease() (entity.Release, error) {
	release, _, err := r.client.Repositories.GetLatestRelease(context.Background(), releaseOwner, releaseRepo)
	if err != nil {
		return entity.Release{}, dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to fetch latest Diflog release", err)
	}
	assets := make([]entity.ReleaseAsset, 0, len(release.Assets))
	for _, asset := range release.Assets {
		assets = append(assets, entity.ReleaseAsset{
			Name:        asset.GetName(),
			DownloadURL: asset.GetBrowserDownloadURL(),
		})
	}
	return entity.Release{
		TagName: release.GetTagName(),
		Assets:  assets,
	}, nil
}
