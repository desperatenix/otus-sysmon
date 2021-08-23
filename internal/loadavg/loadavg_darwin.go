package loadavg

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

type loadavg struct {
	load  [3]uint32
	scale int
}

func get() (*Stats, error) {
	b, err := unix.SysctlRaw("vm.loadavg")
	if err != nil {
		panic(err)
	}
	load := *(*loadavg)(unsafe.Pointer((&b[0])))
	scale := float64(load.scale)
	return &Stats{
		float64(load.load[0]) / scale,
		float64(load.load[1]) / scale,
		float64(load.load[2]) / scale,
	}, nil
}
