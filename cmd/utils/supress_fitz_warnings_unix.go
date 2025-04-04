//go:build !windows
// +build !windows
package utils 

import (
	"os"
	"syscall"
)

func SupressFitzWarnings() {
    devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }

    // Redirect stdout and stderr
    syscall.Dup2(int(devNull.Fd()), int(os.Stderr.Fd()))
}
