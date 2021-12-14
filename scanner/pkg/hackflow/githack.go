package hackflow

import (
	"path/filepath"
)

type gitHack struct {
	baseTool
}

func newGitHack() Tool {
	return &gitHack{
		baseTool: baseTool{
			name:      GIT_HACK,
			desp:      "Git 泄漏利用工具",
			link:      "https://github.com.cnpmjs.org/lijiejie/GitHack",
			installer: GetGit(),
		},
	}
}

func GetGitHack() *gitHack {
	return container.Get(GIT_HACK).(*gitHack)
}

func (g *gitHack) GetExecPath() (string, error) {
	path, err := g.baseTool.GetExecPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, "GitHack.py"), nil
}
