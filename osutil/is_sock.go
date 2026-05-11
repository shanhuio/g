package osutil

import (
	"os"
)

// IsSock checks if a file a unix domain socket.
func IsSock(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return stat.Mode()&os.ModeSocket != 0, nil
}
