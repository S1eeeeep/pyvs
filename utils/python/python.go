package python

import (
	"fmt"
	"os"

	"github.com/S1eeeeep/pyvs/utils/file"
)

func GetInstalled(root string) []string {
	list := make([]string, 0)
	files, _ := os.ReadDir(root)
	for i := len(files) - 1; i >= 0; i-- {
		if files[i].IsDir() {
			list = append(list, files[i].Name())
		}
	}
	return list
}

func IsVersionInstalled(root string, version string) bool {
	isInstalled := file.Exists(fmt.Sprintf("%s\\%s\\python.exe", root, version))
	return isInstalled
}
