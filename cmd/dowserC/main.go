package main

import (
	"os"
	"path/filepath"

	"github.com/Akinoyuu/DowserCE/pkg/launcher"
	"github.com/caarlos0/log"
)

const version = "v1.0.2"

func main() {
	log.Infof("DowserCE 版本: %s", version)
	log.Infof("运行路径：%s\n", launcher.CWD)

	err := launcher.CheckName(launcher.CWD)
	if err != nil {
		log.Error("路径中含有非ASCII字符, 会导致模组无法加载")
		log.Error("即将启动游戏, 3秒后自动退出...")
		launcher.RunDowser()
		return
	}

	launcher.RenameInvalidModFolder(launcher.ModDir)
	entries, err := os.ReadDir(launcher.ModDir)
	if err != nil {
		log.WithError(err).Fatal("获取Mod失败, 可能是因为权限不足")
	}

	launcher.DeleteDotModFiles(filepath.Join(launcher.DataDir, "mod"))
	modNames := []string{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modName, err := launcher.GenerateDotModFile(filepath.Join(launcher.ModDir, entry.Name()))
		if err != nil {
			continue
		}

		modNames = append(modNames, modName)
	}

	if len(modNames) != 0 {
		log.Infof("共计加载: %d个Mod", len(modNames))
		log.Info("他们分别是:")
		log.IncreasePadding()

		for _, modName := range modNames {
			log.Infof("%s", modName)
		}

		log.ResetPadding()
	}

	log.Info("即将启动游戏, 3秒后自动退出...")
	launcher.RunDowser()
}
