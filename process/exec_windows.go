// +build windows

package process

import (
	"github.com/pharosnet/dalg/logger"
	"os"
	"os/exec"
	"strings"
)

func formatAndVet(dir string) (err error) {
	pkgName := ""
	if goPaths, has := os.LookupEnv("GOPATH"); has {
		paths := strings.SplitN(goPaths, ":", -1)
		for _, path := range paths {
			path = strings.TrimSpace(path) + "/src/"
			if strings.Index(dir, path) > -1 {
				pkgName = strings.Replace(dir, path, "", 1)
				break
			}
		}
	}
	if pkgName == "" {
		if goPaths, has := os.LookupEnv("GOROOT"); has {
			paths := strings.SplitN(goPaths, ":", -1)
			for _, path := range paths {
				path = strings.TrimSpace(path) + "/src/"
				if strings.Index(dir, path) > -1 {
					pkgName = strings.Replace(dir, path, "", 1)
					break
				}
			}
		}
	}
	vetCmd := exec.Command("go.exe", "vet", "-v", pkgName)
	if err = vetCmd.Run(); err != nil {
		if bb, oErr := vetCmd.Output(); oErr == nil {
			logger.Log().Println(string(bb))
		}
		return
	}
	if bb, oErr := vetCmd.Output(); oErr == nil {
		logger.Log().Println(string(bb))
	}

	fmtCmd := exec.Command("gofmt.exe", "-w", dir)
	if err = fmtCmd.Run(); err != nil {
		if bb, oErr := fmtCmd.Output(); oErr == nil {
			logger.Log().Println(string(bb))
		}
		return
	}
	if bb, oErr := fmtCmd.Output(); oErr == nil {
		logger.Log().Println(string(bb))
	}
	return
}
