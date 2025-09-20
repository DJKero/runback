package fs

import (
	"os"
	"strings"

	"runback/error"
)

func ReadFileWhole(f string) string {
	data, err := os.ReadFile(f)
	error.ErrorLog(err)
	dataString := string(data)
	return strings.TrimSuffix(dataString, "\n")
}
