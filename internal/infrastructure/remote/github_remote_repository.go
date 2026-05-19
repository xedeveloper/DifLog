package remote

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v66/github"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	dlerrors "github.com/xedeveloper/DifLog/pkg/errors"
	"golang.org/x/oauth2"
)

type GitHubRemoteRepository struct {
	token string
}

func NewGitHubRemoteRepository(token string) *GitHubRemoteRepository {
	return &GitHubRemoteRepository{token: token}
}

func (r *GitHubRemoteRepository) newClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: r.token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func parseRemoteURL(remoteURL string) (owner, repo string, err error) {
	remoteURL = strings.TrimSuffix(remoteURL, ".git")
	if strings.Contains(remoteURL, "@") && strings.Contains(remoteURL, ":") && !strings.Contains(remoteURL, "//") {
		parts := strings.SplitN(remoteURL, ":", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid remote URL: %s", remoteURL)
		}
		pathParts := strings.Split(strings.Trim(parts[1], "/"), "/")
		if len(pathParts) < 2 {
			return "", "", fmt.Errorf("invalid remote URL: %s", remoteURL)
		}
		return pathParts[len(pathParts)-2], pathParts[len(pathParts)-1], nil
	}
	parsed, err := url.Parse(remoteURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid remote URL: %s", remoteURL)
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid remote URL: %s", remoteURL)
	}
	return parts[len(parts)-2], parts[len(parts)-1], nil
}

func (r *GitHubRemoteRepository) PushCommits(remote entity.Remote, commits []entity.Commit, objects map[string][]byte) error {
	ctx := context.Background()
	client := r.newClient(ctx)

	owner, repo, err := parseRemoteURL(remote.URL)
	if err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "invalid remote URL", err)
	}
	branch := "main"
	if len(commits) > 0 && commits[0].Branch != "" {
		branch = commits[0].Branch
	}

	payload := map[string]interface{}{
		"commits": commits,
		"objects": encodeObjects(objects),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to marshal push payload", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	filePath := fmt.Sprintf("diflog/branches/%s/data.json", branch)
	message := "DifLog: push context history"

	existingFile, _, _, _ := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)

	fileContent := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: []byte(encoded),
	}

	if existingFile != nil {
		fileContent.SHA = existingFile.SHA
		_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, filePath, fileContent)
	} else {
		_, _, err = client.Repositories.CreateFile(ctx, owner, repo, filePath, fileContent)
	}

	if err != nil {
		return dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to push to GitHub", err)
	}

	return nil
}

func (r *GitHubRemoteRepository) PullCommits(remote entity.Remote, branch string) ([]entity.Commit, map[string][]byte, error) {
	ctx := context.Background()
	client := r.newClient(ctx)

	owner, repo, err := parseRemoteURL(remote.URL)
	if err != nil {
		return nil, nil, dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "invalid remote URL", err)
	}

	if branch == "" {
		branch = "main"
	}
	filePath := fmt.Sprintf("diflog/branches/%s/data.json", branch)
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
	if err != nil {
		return nil, nil, dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to pull from GitHub", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, nil, dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to decode remote content", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		decoded = []byte(content)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, nil, dlerrors.NewDifLogError(dlerrors.ErrRemoteFailure, "failed to parse remote payload", err)
	}

	commitsData, _ := json.Marshal(payload["commits"])
	var commits []entity.Commit
	json.Unmarshal(commitsData, &commits)

	objectsEncoded, _ := payload["objects"].(map[string]interface{})
	objects := decodeObjects(objectsEncoded)

	return commits, objects, nil
}

func encodeObjects(objects map[string][]byte) map[string]string {
	encoded := make(map[string]string)
	for k, v := range objects {
		encoded[k] = base64.StdEncoding.EncodeToString(v)
	}
	return encoded
}

func decodeObjects(encoded map[string]interface{}) map[string][]byte {
	decoded := make(map[string][]byte)
	for k, v := range encoded {
		if s, ok := v.(string); ok {
			if data, err := base64.StdEncoding.DecodeString(s); err == nil {
				decoded[k] = data
			}
		}
	}
	return decoded
}
