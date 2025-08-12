package launcher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/caarlos0/log"
)

var ErrNoNameFound = errors.New("没有找到Mod的名称")

func GetDotModFileValue(path, key string) (string, error) {
	byteData, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(fmt.Sprintf(`%s\s*=\s*"([^"]+)"`, key))
	match := re.FindStringSubmatch(string(byteData))
	if match == nil {
		return "", ErrNoNameFound
	}

	return match[1], nil
}

func GenerateDotModFile(path string) (string, error) {
	byteData, err := os.ReadFile(filepath.Join(path, "descriptor.mod"))
	if err != nil {
		return "", err
	}

	modName, err := GetDotModFileValue(filepath.Join(path, "descriptor.mod"), "name")
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`[<>:"/\|?*]+`)
	changeName := re.ReplaceAllString(modName, "_")
	dotModPath := filepath.Join(DataDir, "mod", fmt.Sprintf("%s.mod", changeName))

	if _, err := os.Stat(dotModPath); os.IsNotExist(err) {
		file, err := os.Create(dotModPath)
		if err != nil {
			return "", err
		}

		_, err = file.Write(byteData)
		if err != nil {
			return "", err
		}

		_, err = file.WriteString(fmt.Sprintf("\npath=\"%s\"", filepath.ToSlash(path)))
		if err != nil {
			return "", err
		}
		file.Close()
	} else if err != nil {
		log.WithError(err).Fatal("无法读取mod文件")
	}

	return modName, nil
}

func DeleteDotModFiles(path string) error {
	re := regexp.MustCompile(`.*\.mod`)
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		if re.MatchString(fileName) {
			os.Remove(filepath.Join(path, fileName))
		}
	}

	return nil
}
