package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
	"github.com/xedeveloper/DifLog/internal/infrastructure/ai"
	"github.com/xedeveloper/DifLog/internal/infrastructure/diffing"
	"github.com/xedeveloper/DifLog/internal/infrastructure/hashing"
	"github.com/xedeveloper/DifLog/internal/infrastructure/remote"
	"github.com/xedeveloper/DifLog/internal/infrastructure/storage"
	addUC "github.com/xedeveloper/DifLog/internal/usecase/add"
	branchUC "github.com/xedeveloper/DifLog/internal/usecase/branch"
	checkoutUC "github.com/xedeveloper/DifLog/internal/usecase/checkout"
	commitUC "github.com/xedeveloper/DifLog/internal/usecase/commit"
	detectUC "github.com/xedeveloper/DifLog/internal/usecase/detect"
	diffUC "github.com/xedeveloper/DifLog/internal/usecase/diff"
	initUC "github.com/xedeveloper/DifLog/internal/usecase/initialize"
	logUC "github.com/xedeveloper/DifLog/internal/usecase/log"
	pullUC "github.com/xedeveloper/DifLog/internal/usecase/pull"
	pushUC "github.com/xedeveloper/DifLog/internal/usecase/push"
	skillUC "github.com/xedeveloper/DifLog/internal/usecase/skill"
	statusUC "github.com/xedeveloper/DifLog/internal/usecase/status"
)

type Container struct {
	ProjectRoot string
	Initializer *storage.LocalStorageInitializer
	DifLogDir   string

	ContextRepo      *storage.ContextLocalRepository
	CommitRepo       *storage.CommitLocalRepository
	BranchRepo       *storage.BranchLocalRepository
	StagingRepo      *storage.StagingLocalRepository
	SkillContextRepo *storage.SkillContextLocalRepository

	HasherService  *hashing.SHA256HasherService
	DifferService  *diffing.TextDifferService
	AIDetector     *ai.AIToolDetector
	ClaudeAdapter  *ai.ClaudeCodeAdapter
	OpenCodeAdapter *ai.OpenCodeAdapter
	CopilotAdapter *ai.CopilotCLIAdapter

	RemoteRepo *remote.GitHubRemoteRepository

	InitUC          *initUC.InitializeUseCase
	AddUC           *addUC.AddContextUseCase
	CommitUC        *commitUC.CommitContextUseCase
	LogUC           *logUC.ViewCommitLogUseCase
	DiffUC          *diffUC.DiffCommitsUseCase
	CheckoutUC      *checkoutUC.CheckoutCommitUseCase
	StatusUC        *statusUC.ViewStatusUseCase
	DetectUC        *detectUC.DetectAIToolUseCase
	SkillUC         *skillUC.CreateSkillUseCase
	SaveSkillCtxUC  *skillUC.SaveSkillContextUseCase
	PushUC          *pushUC.PushToRemoteUseCase
	PullUC          *pullUC.PullFromRemoteUseCase

	CreateBranchUC *branchUC.CreateBranchUseCase
	ListBranchesUC *branchUC.ListBranchesUseCase
	SwitchBranchUC *branchUC.SwitchBranchUseCase
}

func NewContainer(projectRoot string) *Container {
	initializer := storage.NewLocalStorageInitializer(projectRoot)
	difLogDir := initializer.DifLogDir()

	contextRepo := storage.NewContextLocalRepository(difLogDir)
	commitRepo := storage.NewCommitLocalRepository(difLogDir)
	branchRepo := storage.NewBranchLocalRepository(difLogDir, initializer)
	stagingRepo := storage.NewStagingLocalRepository(difLogDir)
	skillContextRepo := storage.NewSkillContextLocalRepository(difLogDir)

	hasherSvc := hashing.NewSHA256HasherService()
	differSvc := diffing.NewTextDifferService()
	aiDetector := ai.NewAIToolDetector()
	claudeAdapter := ai.NewClaudeCodeAdapter()
	openCodeAdapter := ai.NewOpenCodeAdapter()
	copilotAdapter := ai.NewCopilotCLIAdapter()

	githubToken := os.Getenv("DIFLOG_GITHUB_TOKEN")
	remoteRepo := remote.NewGitHubRemoteRepository(githubToken)

	c := &Container{
		ProjectRoot:      projectRoot,
		Initializer:      initializer,
		DifLogDir:        difLogDir,
		ContextRepo:      contextRepo,
		CommitRepo:       commitRepo,
		BranchRepo:       branchRepo,
		StagingRepo:      stagingRepo,
		SkillContextRepo: skillContextRepo,
		HasherService:    hasherSvc,
		DifferService:    differSvc,
		AIDetector:       aiDetector,
		ClaudeAdapter:    claudeAdapter,
		OpenCodeAdapter:  openCodeAdapter,
		CopilotAdapter:   copilotAdapter,
		RemoteRepo:       remoteRepo,
	}

	c.InitUC = initUC.NewInitializeUseCase(initializer)
	c.AddUC = addUC.NewAddContextUseCase(contextRepo, stagingRepo, hasherSvc, aiDetector, projectRoot)
	c.CommitUC = commitUC.NewCommitContextUseCase(commitRepo, stagingRepo, contextRepo, branchRepo, hasherSvc)
	c.LogUC = logUC.NewViewCommitLogUseCase(commitRepo, branchRepo)
	c.DiffUC = diffUC.NewDiffCommitsUseCase(commitRepo, differSvc)
	c.CheckoutUC = checkoutUC.NewCheckoutCommitUseCase(commitRepo, projectRoot)
	c.StatusUC = statusUC.NewViewStatusUseCase(stagingRepo, branchRepo)
	c.DetectUC = detectUC.NewDetectAIToolUseCase(aiDetector, projectRoot)
	c.SkillUC = skillUC.NewCreateSkillUseCase(claudeAdapter, openCodeAdapter, copilotAdapter, projectRoot)
	c.SaveSkillCtxUC = skillUC.NewSaveSkillContextUseCase(skillContextRepo, hasherSvc)
	c.PushUC = pushUC.NewPushToRemoteUseCase(commitRepo, contextRepo, branchRepo, remoteRepo, initializer)
	c.PullUC = pullUC.NewPullFromRemoteUseCase(commitRepo, contextRepo, branchRepo, remoteRepo, initializer)
	c.CreateBranchUC = branchUC.NewCreateBranchUseCase(branchRepo, commitRepo)
	c.ListBranchesUC = branchUC.NewListBranchesUseCase(branchRepo)
	c.SwitchBranchUC = branchUC.NewSwitchBranchUseCase(branchRepo)

	return c
}

func (c *Container) RequireInitialized() error {
	if !c.Initializer.IsInitialized() {
		return fmt.Errorf("DifLog is not initialized. Run 'diflog init' first")
	}
	return nil
}

func (c *Container) LoadAITool() entity.AITool {
	config, err := c.Initializer.LoadConfig()
	if err != nil {
		return entity.AIToolUnknown
	}
	return config.AITool
}

func NewRootCommand() *cobra.Command {
	projectRoot, err := os.Getwd()
	if err != nil {
		projectRoot = "."
	}

	container := NewContainer(projectRoot)

	rootCmd := &cobra.Command{
		Use:   "diflog",
		Short: "DifLog — AI Context Version Control",
		Long: `DifLog is a Git-like version control tool for AI contexts.
Track, commit, and restore your ClaudeCode and GitHub Copilot CLI context files.`,
	}

	rootCmd.AddCommand(
		newInitCommand(container),
		newAddCommand(container),
		newCommitCommand(container),
		newLogCommand(container),
		newDiffCommand(container),
		newCheckoutCommand(container),
		newBranchCommand(container),
		newPushCommand(container),
		newPullCommand(container),
		newSkillCommand(container),
		newDetectCommand(container),
		newStatusCommand(container),
	)

	return rootCmd
}
