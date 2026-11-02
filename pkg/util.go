package pkg

import "os/user"

// ExpandHomePath resolves ~ to the current user's home directory when possible.
func ExpandHomePath(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}
	usr, err := user.Current()
	if err != nil || usr == nil {
		return path
	}
	return usr.HomeDir + path[1:]
}
