// +build darwin netbsd freebsd openbsd dragonfly linux

package process

import (
	"github.com/pharosnet/dalg/logger"
	"os/exec"
)

func formatAndVet(dir string) error {
	fmtCmd := exec.Command("gofmt", "-w", dir + "/*.go")
	if bb, err := fmtCmd.Output(); err != nil {
		logger.Log().Println(string(bb))
		return err
	}
	vetCmd := exec.Command("go", "vet", "-v", dir)
	if bb, err := vetCmd.Output(); err != nil {
		logger.Log().Println(string(bb))
		return err
	}
	return nil
}
