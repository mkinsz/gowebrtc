package server

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// Generate CSV with two columns of peerConnectionCount and cpuUsage
func serveReport() {
	file, err := os.OpenFile("report.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	if _, err := file.WriteString("peerConnectionCount, cpuUsage\n"); err != nil {
		panic(err)
	}

	for range time.NewTicker(3 * time.Second).C {
		usage, err := cpu.Percent(0, false)
		if err != nil {
			panic(err)
		} else if len(usage) != 1 {
			panic(fmt.Sprintf("CPU Usage results should have 1 sample, have %d", len(usage)))
		}

		if _, err = file.WriteString(fmt.Sprintf("%d, %f\n", atomic.LoadInt64(&peerConnectionCount), usage[0])); err != nil {
			panic(err)
		}
	}
}
