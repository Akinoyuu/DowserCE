package launcher

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

var ErrNonASCIIChar = errors.New("路径中包含非ASCII字符")

func CheckName(path string) error {
	re := regexp.MustCompile(`[^\x20-\x7E]`)
	match := re.FindStringSubmatch(path)
	if match != nil {
		return ErrNonASCIIChar
	}

	return nil
}

func RenameInvalidModFolder(path string) error {
	entries, err := os.ReadDir(ModDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		re := regexp.MustCompile(`[^\x20-\x7E]`)
		match := re.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}

		name := filepath.Join(path, entry.Name())
		value, err := GetDotModFileValue(filepath.Join(name, "descriptor.mod"), "remote_file_id")
		if err != nil {
			value = re.ReplaceAllString(entry.Name(), "_")
		}

		err = os.Rename(name, filepath.Join(path, value))
		if err != nil {
			return err
		}
	}

	return nil
}
