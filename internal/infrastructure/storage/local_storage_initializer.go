package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/xedeveloper/DifLog/internal/domain/entity"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
)

const DifLogDirName = ".difLog"

type DifLogConfig struct {
	Branch string         `json:"branch"`
	Remote string         `json:"remote"`
	AITool entity.AITool  `json:"aiTool"`
}

type LocalStorageInitializer struct {
	projectRoot string
}

func NewLocalStorageInitializer(projectRoot string) *LocalStorageInitializer {
	return &LocalStorageInitializer{projectRoot: projectRoot}
}

func (i *LocalStorageInitializer) DifLogDir() string {
	return filepath.Join(i.projectRoot, DifLogDirName)
}

func (i *LocalStorageInitializer) IsInitialized() bool {
	_, err := os.Stat(i.DifLogDir())
	return err == nil
}

func (i *LocalStorageInitializer) Initialize(aiTool entity.AITool) error {
	if i.IsInitialized() {
		return dlerrors.NewDifLogError(dlerrors.ErrAlreadyInitialized, "DifLog is already initialized in this directory", nil)
	}

	dirs := []string{
		i.DifLogDir(),
		filepath.Join(i.DifLogDir(), "staging"),
		filepath.Join(i.DifLogDir(), "objects"),
		filepath.Join(i.DifLogDir(), "refs", "heads"),
		filepath.Join(i.DifLogDir(), "refs", "remote", "origin"),
		filepath.Join(i.DifLogDir(), "logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to create directory: "+dir, err)
		}
	}

	if err := i.writeHEAD("main"); err != nil {
		return err
	}

	if err := i.writeConfig(DifLogConfig{Branch: "main", Remote: "", AITool: aiTool}); err != nil {
		return err
	}

	mainRefPath := filepath.Join(i.DifLogDir(), "refs", "heads", "main")
	if err := os.WriteFile(mainRefPath, []byte(""), 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to create main ref", err)
	}

	stagingIndex := entity.StagingIndex{Entries: []entity.StagedContext{}}
	indexData, _ := json.Marshal(stagingIndex)
	indexPath := filepath.Join(i.DifLogDir(), "staging", "index.json")
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to create staging index", err)
	}

	commitsPath := filepath.Join(i.DifLogDir(), "logs", "commits.json")
	if err := os.WriteFile(commitsPath, []byte("[]"), 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to create commits log", err)
	}

	return nil
}

func (i *LocalStorageInitializer) writeHEAD(branch string) error {
	headPath := filepath.Join(i.DifLogDir(), "HEAD")
	content := "ref: refs/heads/" + branch
	if err := os.WriteFile(headPath, []byte(content), 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to write HEAD", err)
	}
	return nil
}

func (i *LocalStorageInitializer) writeConfig(config DifLogConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to marshal config", err)
	}
	configPath := filepath.Join(i.DifLogDir(), "config.json")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to write config", err)
	}
	return nil
}

func (i *LocalStorageInitializer) LoadConfig() (DifLogConfig, error) {
	configPath := filepath.Join(i.DifLogDir(), "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DifLogConfig{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to read config", err)
	}
	var config DifLogConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return DifLogConfig{}, dlerrors.NewDifLogError(dlerrors.ErrStorageFailure, "failed to parse config", err)
	}
	return config, nil
}

func (i *LocalStorageInitializer) SaveConfig(config DifLogConfig) error {
	return i.writeConfig(config)
}

func (i *LocalStorageInitializer) LoadRemote() (string, error) {
	config, err := i.LoadConfig()
	if err != nil {
		return "", err
	}
	return config.Remote, nil
}

func (i *LocalStorageInitializer) SaveRemote(remote string) error {
	config, err := i.LoadConfig()
	if err != nil {
		return err
	}
	config.Remote = remote
	return i.SaveConfig(config)
}

func (i *LocalStorageInitializer) UpdateHEAD(branch string) error {
	return i.writeHEAD(branch)
}
