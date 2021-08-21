package cpu

import (
	"fmt"
	"unsafe"
)

// #include <mach/mach_host.h>
// #include <mach/host_info.h>
import "C"

// Get cpu statistics
func get() (*Stats, error) {
	return collectCPUStats()
}


func collectCPUStats() (*Stats, error) {
	var cpuLoad C.host_cpu_load_info_data_t
	var count C.mach_msg_type_number_t = C.HOST_CPU_LOAD_INFO_COUNT
	ret := C.host_statistics(C.host_t(C.mach_host_self()), C.HOST_CPU_LOAD_INFO, C.host_info_t(unsafe.Pointer(&cpuLoad)), &count)
	if ret != C.KERN_SUCCESS {
		return nil, fmt.Errorf("host_statistics failed: %d", ret)
	}



	cpu := Stats{}

	user := float64(cpuLoad.cpu_ticks[C.CPU_STATE_USER])
	system := float64(cpuLoad.cpu_ticks[C.CPU_STATE_SYSTEM])
	idle := float64(cpuLoad.cpu_ticks[C.CPU_STATE_IDLE])

	total := user+system+idle

	cpu.User = user/total*100
	cpu.System = system/total*100
	cpu.Idle = idle/total*100
	cpu.Total = total

	return &cpu, nil
}