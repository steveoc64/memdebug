package memdebug

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"time"
)

var lastMemm uint64
var lastMemmm uint64
var printMutex sync.Mutex

const (
	mmOff     = "[0m"
	mmBlue    = "[36m"
	mmOrange  = "[33m"
	mmYellow  = "[1;33m"
	mmGreen   = "[32m"
	mmMagenta = "[35m"
	mmRed     = "[1;31m"
)

var isProfiling bool

func Profile() {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("cpu.pprof", err)
	}
	pprof.StartCPUProfile(f)
	isProfiling = true
}

func Print(t time.Time, what ...interface{}) {
	printMutex.Lock()
	defer printMutex.Unlock()

	// OTT memory hacks
	ms1 := &runtime.MemStats{}
	ms2 := &runtime.MemStats{}
	if !isProfiling {
		runtime.ReadMemStats(ms1)
		runtime.GC()
		debug.FreeOSMemory()
	}
	runtime.ReadMemStats(ms2)
	mmV := ms2.Alloc - lastMemm
	cmm := mmYellow
	cmmm := mmOff
	if ms2.Alloc < lastMemm {
		mmV = lastMemm - ms2.Alloc
		cmm = mmGreen
	}
	if ms2.Sys > lastMemmm {
		cmmm = mmRed
	}

	// build up a string and print it once, otherwise the output from different
	// threads can easily get gemogrified together
	s := fmt.Sprintf("%s%12s%s [%s%4d%s] (%s%8v%s):%10v:%s%10v%s:%10v <- %s  %+v%s",
		//mmBlue, mtype.method.Name, mmOff,
		mmOrange, time.Since(t), mmOff,
		mmBlue, runtime.NumGoroutine(), mmOff,
		cmm, mmV, mmOff,
		ms2.Alloc,
		cmmm, ms2.Sys, mmOff,
		ms2.StackInuse,
		mmMagenta,
		what,
		mmOff)

	fmt.Println(s)

	lastMemm = ms2.Alloc
	lastMemmm = ms2.Sys

}

var pversion = 1

func DumpProfile() {
	pprof.StopCPUProfile()
	cmd := exec.Command("go", "tool", "pprof", "-pdf", "./cpu.pprof")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(fmt.Sprintf("cpu-%03d.pdf", pversion))
	if err == nil {
		io.Copy(f, stdout)
		f.Close()
	}
	cmd.Wait()
	pversion++
	pprof.StartCPUProfile(f)
}

func WriteProfile() {
	pprof.StopCPUProfile()
	cmd := exec.Command("go", "tool", "pprof", "-pdf", "./cpu.pprof")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("cpu.pdf")
	if err == nil {
		io.Copy(f, stdout)
		f.Close()
	}
	cmd.Wait()
}
