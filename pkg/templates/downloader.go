package templates

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func cloneTemplatesRepo(repoPath string, force bool) error {
	repo, err := git.PlainOpen(repoPath)
	if err == nil {
		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}

		status, err := worktree.Status()
		if err != nil {
			return err
		}

		if !force {
			if !status.IsClean() {
				return fmt.Errorf("deteceted uncomitted changes in %s", repoPath)
			}
		}

		err = worktree.Pull(&git.PullOptions{
			RemoteName: "origin",
			Force:      true,
			Progress:   os.Stdout,
		})

		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}

		return nil
	}

	if err := os.MkdirAll(repoPath, 0750); err != nil {
		return err
	}

	_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      "https://github.com/HappyHackingSpace/vt-templates",
		Progress: os.Stdout,
		Depth:    1,
	})

	if err != nil {
		return err
	}

	return nil
}
