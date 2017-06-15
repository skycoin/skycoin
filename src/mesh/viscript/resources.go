package viscript

import (
	"runtime"
)

func (self *ViscriptServer) GetResources() (float64, uint64, error) {
	return 0, getMemStats(), nil
}

func getMemStats() uint64 {
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	return ms.Alloc
}

func getCPUProfile() {

}
