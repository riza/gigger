package gigger

const (
	//VERSION current version number
	VERSION = "1.0"
)

var (
	//GitFolderStructure Architecture and content of the Git folder
	GitFolderStructure = map[string]interface{}{
		".git/COMMIT_EDITMSG": true,
		".git/HEAD":           true,
		".git/config":         true,
		".git/description":    true,
		".git/index":          true,
		".git/hooks/": map[string]bool{
			"applypatch-msg.sample":     true,
			"commit-msg.sample":         true,
			"fsmonitor-watchman.sample": true,
			"post-update.sample":        true,
			"pre-applypatch.sample":     true,
			"pre-commit.sample":         true,
			"pre-merge-commit.sample":   true,
			"pre-push.sample":           true,
			"pre-rebase.sample":         true,
			"pre-receive.sample":        true,
			"prepare-commit-msg.sample": true,
			"update.sample":             true,
		},
		".git/info/": map[string]bool{
			"exclude": true,
		},
		".git/logs/": map[string]bool{
			"HEAD":                                   true,
			"refs/heads":                             false,
			"refs/heads/master":                      true,
			"refs/heads/remotes":                     false,
			"refs/heads/remotes/heads":               false,
			"refs/heads/remotes/heads/origin":        false,
			"refs/heads/remotes/heads/origin/master": true,
		},
		".git/objects/": false,
		".git/refs/": map[string]bool{
			"heads":                 false,
			"heads/master":          true,
			"remotes":               false,
			"remotes/origin":        false,
			"remotes/origin/master": true,
			"tags":                  false,
		},
	}
)
