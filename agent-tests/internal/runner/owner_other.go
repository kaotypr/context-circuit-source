//go:build !unix

package runner

import "os"

func ownedByCurrentUser(os.FileInfo) bool { return false }
