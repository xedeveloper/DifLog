package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xedeveloper/DifLog/internal/domain/entity"
)

func newSkillCommand(c *Container) *cobra.Command {
	skillCmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage AI tool skills and context snapshots",
	}

	skillCmd.AddCommand(
		newSkillCreateCommand(c),
		newSkillSaveCommand(c),
	)

	return skillCmd
}

func newSkillCreateCommand(c *Container) *cobra.Command {
	return &cobra.Command{
		Use:   "create [ai-tool]",
		Short: "Create a skill for the detected or specified AI tool",
		Long: `Create a DifLog skill/instructions file for the configured AI tool.

For ClaudeCode:        creates ~/.claude/skills/diflog/SKILL.md
For OpenCode:          creates ~/.config/opencode/skills/diflog/SKILL.md
For GitHub Copilot CLI: appends to .github/copilot-instructions.md

The skill includes instructions for the AI to:
  - Read .difLog/logs/commits.json to present context history
  - Save chat context snapshots to .difLog/skill-contexts/{hash}.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			var aiTool entity.AITool
			if len(args) > 0 {
				aiTool = entity.AITool(args[0])
			} else {
				aiTool = c.LoadAITool()
			}

			if aiTool == entity.AIToolUnknown || aiTool == "" {
				return fmt.Errorf("could not determine AI tool. Specify: diflog skill create [claude|opencode|copilot]")
			}

			result, err := c.SkillUC.Execute(aiTool)
			if err != nil {
				return err
			}

			fmt.Printf("✓ %s\n", result.Message)
			return nil
		},
	}
}

func newSkillSaveCommand(c *Container) *cobra.Command {
	var contentFlag string
	var sourceFlag string

	cmd := &cobra.Command{
		Use:   "save [file]",
		Short: "Save a chat context snapshot to .difLog/skill-contexts/",
		Long: `Snapshot the current AI chat context into .difLog/skill-contexts/{hash}.json.

The snapshot stores: hash (SHA-256), timestamp, content, and source label.
This allows the AI skill to reference past chat contexts by hash.

Examples:
  diflog skill save context.md
  diflog skill save --content "Summary of session..." --source "chat-2024-01-15"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.RequireInitialized(); err != nil {
				return err
			}

			var result interface{ GetHash() string }

			if len(args) > 0 {
				r, err := c.SaveSkillCtxUC.ExecuteFromFile(args[0])
				if err != nil {
					return err
				}
				fmt.Printf("✓ Context snapshot saved\n")
				fmt.Printf("  Hash:     %s\n", r.Hash)
				fmt.Printf("  Location: %s\n", r.FilePath)
				_ = result
				return nil
			}

			if contentFlag == "" {
				return fmt.Errorf("provide a file path or use --content to supply context inline")
			}

			source := sourceFlag
			if source == "" {
				source = "manual"
			}

			r, err := c.SaveSkillCtxUC.ExecuteFromContent(contentFlag, source)
			if err != nil {
				return err
			}

			fmt.Printf("✓ Context snapshot saved\n")
			fmt.Printf("  Hash:     %s\n", r.Hash)
			fmt.Printf("  Location: %s\n", r.FilePath)
			return nil
		},
	}

	cmd.Flags().StringVar(&contentFlag, "content", "", "Inline context content to snapshot")
	cmd.Flags().StringVar(&sourceFlag, "source", "", "Label for the context source (default: manual)")
	return cmd
}
