package launcher

import (
	"os/exec"
	"path/filepath"
	"time"

	"github.com/caarlos0/log"
)

func RunDowser() {
	err := exec.Command(filepath.Join(CWD, "dowser.exe")).Start()
	if err != nil {
		log.WithError(err).Fatal("启动客户端失败, 可能是因为权限不足")
	}

	time.Sleep(3 * time.Second)
}
