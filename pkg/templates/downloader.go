package templates

import (
	"errors"
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

		if !force && !status.IsClean() {
			return fmt.Errorf("detected uncommitted changes in %s", repoPath)
		}

		err = worktree.Pull(&git.PullOptions{
			RemoteName: "origin",
			Force:      true,
		})

		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}

		return nil
	}

	err = os.RemoveAll(repoPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(repoPath, 0750); err != nil {
		return err
	}

	_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:   TemplateRemoteRepoistory,
		Depth: 1,
	})

	if err != nil {
		removeErr := os.RemoveAll(repoPath)
		if removeErr != nil {
			return errors.Join(err, fmt.Errorf("cleanup failed: %w", removeErr))
		}
		return err
	}

	return nil
}
