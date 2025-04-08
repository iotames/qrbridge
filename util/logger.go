package util

import (
	"path/filepath"
	"sync"

	"github.com/iotames/miniutils"
	"github.com/iotames/qrbridge/conf"
)

var (
	once sync.Once
	lg   *miniutils.Logger
)

func GetLogger() *miniutils.Logger {
	once.Do(func() {
		if lg == nil {
			lg = miniutils.NewLogger(filepath.Join(conf.RuntimeDir, "logs"))
		}
	})
	return lg
}
