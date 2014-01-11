package main

import (
	"fmt"
	"github.com/op/go-logging"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"syscall"
)

import (
	"./src/cli/"
	// "./src/coin/"
	"./src/daemon/"
	"./src/examples/"
	"./src/gui/"
)

var (
	logger     = logging.MustGetLogger("skycoin.main")
	logFormat  = "[%{module}:%{level}] %{message}"
	logModules = []string{
		"skycoin.main",
		"skycoin.daemon",
		"skycoin.coin",
		"skycoin.gui",
		"skycoin.util",
		"gnet",
		"pex",
	}
)

func printProgramStatus() {
	fmt.Println("Program Status:")
	debug.PrintStack()
	p := pprof.Lookup("goroutine")
	f, err := os.Create("goroutine.prof")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = p.WriteTo(f, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func catchInterrupt(quit chan int) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	shutdown(quit)
}

// Catches SIGUSR1 and prints internal program state
func catchDebug() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGUSR1)
	for {
		select {
		case <-sigchan:
			printProgramStatus()
		}
	}
}

func shutdown(quit chan int) {
	logger.Info("Shutting down\n")
	daemon.Shutdown(cli.DataDirectory)
	logger.Info("Goodbye\n")
	quit <- 1
}

// func initSettings() {
//     sb.InitSettings()
//     sb.Settings.Load()
//     we resave the settings, in case they were not found and had to be generated
//     sb.Settings.Save()
// }

func initLogging(level logging.Level) {
	format := logging.MustStringFormatter(logFormat)
	logging.SetFormatter(format)
	for _, s := range logModules {
		logging.SetLevel(level, s)
	}
}

func initProfiling() {
	if cli.ProfileCPU {
		f, err := os.Create(cli.ProfileCPUFile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if cli.HTTPProf {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
}

func run_tests() {
	sb_examples.Run()
}

func main() {

	if true {
		run_tests()
	}

	cli.ParseArgs()
	initProfiling()
	initLogging(cli.LogLevel)

	// If the user Ctrl-C's, shutdown properly
	quit := make(chan int)
	go catchInterrupt(quit)
	// Watch for SIGUSR1
	go catchDebug()

	daemon.Init(cli.Port, cli.DataDirectory)

	if cli.ConnectTo != "" {
		_, err := daemon.Pool.Connect(cli.ConnectTo)
		if err != nil {
			log.Panic(err)
		}
	}

	if !cli.DisableGUI {
		go gui.LaunchGUI()
	}

	<-quit
}
