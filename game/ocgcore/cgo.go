//go:build darwin || freebsd || linux || windows

package ocgcore

import (
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

func getSystemLibrary() string {
	switch runtime.GOOS {
	case "darwin":
		return "/usr/lib/libSystem.B.dylib"
	case "linux":
		return "libc.so.6"
	case "freebsd":
		return "libc.so.7"
	case "windows":
		return "ucrtbase.dll"
	default:
		panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
	}
}

func OCGApi() {
	libc, err := openLibrary("E:\\Go\\gopath\\go-ygocore\\libs\\ocgcore.dll")
	if err != nil {
		panic(err)
	}
	var cgo = new(CGO)
	purego.RegisterLibFunc(&cgo.CreateDuel, libc, "create_duel")
	fmt.Println(cgo.CreateDuel(1000))
}

type CGO struct {
	CreateDuel func(seed int32) uintptr
}
