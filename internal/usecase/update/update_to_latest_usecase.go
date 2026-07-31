package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/xedeveloper/DifLog/internal/domain/repository"
	"github.com/xedeveloper/DifLog/internal/version"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

var releaseAssetNames = map[string]string{
	"linux/amd64": "DifLog-linux-amd-64",
	"darwin/amd64": "DifLog-darwin-amd64",
	"darwin/arm64": "DifLog-darwin-arm-64",
}

type UpdateResult struct{
	CurrentVersion string
	LatestVersion string
	Updated bool
	InstalledPath string
}

type UpdateToLatestUseCase struct {
	releaseRepo repository.ReleaseRepository
}

func NewUpdateToLatestUseCase(releaseRepo repository.ReleaseRepository) *UpdateToLatestUseCase{
	return  &UpdateToLatestUseCase{releaseRepo: releaseRepo}
}

func (uc *UpdateToLatestUseCase) Execute() (UpdateResult, error){
	platform := runtime.GOOS + "/" + runtime.GOARCH
	assetName,ok := releaseAssetNames[platform]
	if !ok {
		return UpdateResult{}, dlerrors.NewDifLogError(
			dlerrors.ErrUnsupportedPlatform,
			fmt.Sprintf("unsupported platform %s: diflog update only supported on macOS & linux",platform),
			nil,
		) 
	}
	release, err := uc.releaseRepo.GetLatestRelease()
	if err != nil{
		return  UpdateResult{},err
	}

	current := version.Version
	if current == release.TagName {
		return  UpdateResult{CurrentVersion: current, LatestVersion: release.TagName, Updated: false}, nil
	}
	downloadURL := ""
	for _,asset := range release.Assets{
		if asset.Name == assetName{
			downloadURL = asset.DownloadURL
			break
		}
	}
	if downloadURL == "" {
		return UpdateResult{}, dlerrors.NewDifLogError(
			dlerrors.ErrRemoteFailure,
			fmt.Sprintf("release %s has no asset named %s", release.TagName, assetName),
			nil,
		)
	}

	execPath,err := os.Executable()
	if err != nil{
		return  UpdateResult{},dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to locate the running diflog binary",err)
	}
	execPath,err = filepath.EvalSymlinks(execPath)
	if err != nil{
		return  UpdateResult{}, dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to resolve the running diflog binary path", err)
	}

	if err := downloadAndReplace(downloadURL, execPath); err != nil {
		return  UpdateResult{}, err
	}

	return UpdateResult{
		CurrentVersion:  current,
		LatestVersion: release.TagName,
		Updated: true,
		InstalledPath: execPath,
	}, nil
}

func downloadAndReplace(downloadURL string, execPath string) error{
	resp, err := http.Get(downloadURL)
	if err != nil{
		return  dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to download the latest diflog release", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dlerrors.NewDifLogError(
			dlerrors.ErrRemoteFailure,
			fmt.Sprintf("failed to download the latest diflog release: HTTP %d", resp.StatusCode),
			nil,
		)
	}

	dir := filepath.Dir(execPath)
	tmpFile, err := os.CreateTemp(dir, "diflog-update-*")
	if err != nil{
		return dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to stage the downloaded library",err)
	}

	tmpPath := filepath.Dir(execPath)
	defer os.Remove(tmpPath)

	if _,err := io.Copy(tmpFile, resp.Body); err != nil{
		tmpFile.Close()
		return  dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to write the downloaded binary", err)
	}

	if err := tmpFile.Close(); err != nil{
		return  dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to finalise the downloaded binary", err)
	}
	if err := os.Chmod(tmpPath,0755); err != nil{
		return dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to mark the downloaded binary executable",err)
	}
	if err := os.Rename(tmpPath, execPath); err != nil{
		return  dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to replace the running diflog binary",err)
	}
	return nil
}
