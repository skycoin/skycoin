package skycoin

import (
    "fmt"
    "github.com/op/go-logging"
    "log"
    "net/http"
    _ "net/http/pprof"
    "os"
    "os/signal"
    "runtime/pprof"
    "syscall"
    "time"
)

//TODO, move /src/skycoin to /cmd/skycoin folder

import (
    "github.com/skycoin/skycoin/src/cli"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/daemon"
    "github.com/skycoin/skycoin/src/gui"
    "github.com/skycoin/skycoin/src/visor"
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
        "skycoin.visor",
        "gnet",
        "pex",
    }
)

func printProgramStatus() {
    fn := "goroutine.prof"
    logger.Debug("Writing goroutine profile to %s", fn)
    p := pprof.Lookup("goroutine")
    f, err := os.Create(fn)
    defer f.Close()
    if err != nil {
        logger.Error("%v", err)
        return
    }
    err = p.WriteTo(f, 2)
    if err != nil {
        logger.Error("%v", err)
        return
    }
}

func catchInterrupt(quit chan int) {
    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, os.Interrupt)
    <-sigchan
    signal.Stop(sigchan)
    quit <- 1
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

// func initSettings() {
//     sb.InitSettings()
//     sb.Settings.Load()
//     we resave the settings, in case they were not found and had to be generated
//     sb.Settings.Save()
// }

func initLogging(level logging.Level, color bool) {
    format := logging.MustStringFormatter(logFormat)
    logging.SetFormatter(format)
    for _, s := range logModules {
        logging.SetLevel(level, s)
    }
    stdout := logging.NewLogBackend(os.Stdout, "", 0)
    stdout.Color = color
    logging.SetBackend(stdout)
}

func initProfiling(httpProf, profileCPU bool, profileCPUFile string) {
    if profileCPU {
        f, err := os.Create(profileCPUFile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    if httpProf {
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
}

func configureDaemon(c *cli.Config) *daemon.Config {
    dc := daemon.NewConfig()
    dc.Peers.DataDirectory = c.DataDirectory
    dc.DHT.Disabled = c.DisableDHT
    dc.DHT.Port = c.Port
    dc.Pool.Port = c.Port
    connectionRate := time.Second * time.Duration(c.OutgoingConnectionsRate)
    if c.OutgoingConnectionsRate == 0 {
        connectionRate = time.Millisecond * 30
    }
    dc.Daemon.OutgoingRate = connectionRate
    dc.Visor.Config.IsMaster = c.MasterChain
    dc.Visor.Config.CanSpend = c.CanSpend
    dc.Visor.Config.WalletFile = c.WalletFile
    dc.Visor.Config.WalletSizeMin = c.WalletSizeMin
    dc.Visor.Config.BlockchainFile = c.BlockchainFile
    dc.Visor.Config.BlockSigsFile = c.BlockSigsFile
    dc.Visor.Config.GenesisSignature = coin.MustSigFromHex(c.GenesisSignature)
    dc.Visor.Config.GenesisTimestamp = c.GenesisTimestamp
    if c.MasterChain {
        // The master chain should be reluctant to expire transactions
        dc.Visor.Config.UnconfirmedRefreshRate = time.Hour * 4096
    }

    dc.Visor.MasterKeysFile = c.MasterKeys
    if c.MasterChain {
        err := dc.Visor.LoadMasterKeys()
        if err != nil {
            log.Panicf("Failed to load master keys: %v", err)
        }
    } else {
        w := &visor.ReadableWalletEntry{
            Address: c.GenesisAddress,
            Public:  c.MasterPublic,
        }
        dc.Visor.Config.MasterKeys = visor.WalletEntryFromReadable(w)
    }
    return dc
}

func Run(args cli.Args) {
    c := cli.ParseArgs(args)
    initProfiling(c.HTTPProf, c.ProfileCPU, c.ProfileCPUFile)
    initLogging(c.LogLevel, c.ColorLog)

    // If the user Ctrl-C's, shutdown properly
    quit := make(chan int)
    go catchInterrupt(quit)
    // Watch for SIGUSR1
    go catchDebug()

    dconf := configureDaemon(c)
    d := daemon.NewDaemon(dconf)

    stopDaemon := make(chan int)
    go d.Start(stopDaemon)

    if c.ConnectTo != "" {
        _, err := d.Pool.Pool.Connect(c.ConnectTo)
        if err != nil {
            log.Panic(err)
        }
    }

    if !c.DisableGUI {
        go gui.LaunchGUI(d)
    }

    host := fmt.Sprintf("%s:%d", c.WebInterfaceAddr, c.WebInterfacePort)

    if c.WebInterface {
        if c.WebInterfaceHTTPS {
            // Verify cert/key parameters, and if neither exist, create them
            errs := gui.CreateCertIfNotExists(host, c.WebInterfaceCert,
                c.WebInterfaceKey)
            if len(errs) != 0 {
                for _, err := range errs {
                    logger.Error(err.Error())
                }
            } else {
                go gui.LaunchWebInterfaceHTTPS(host, c.GUIDirectory, d,
                    c.WebInterfaceCert, c.WebInterfaceKey)
            }
        } else {
            go gui.LaunchWebInterface(host, c.GUIDirectory, d)
        }
    }

    <-quit
    stopDaemon <- 1

    logger.Info("Shutting down")
    d.Shutdown()
    logger.Info("Goodbye")
}
