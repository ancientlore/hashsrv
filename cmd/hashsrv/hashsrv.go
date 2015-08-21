package main

import (
	"flag"
	"fmt"
	"github.com/ancientlore/flagcfg"
	"github.com/facebookgo/flagenv"
	"github.com/kardianos/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
)

var cpuprofile string
var memprofile string
var wsHashSrv service.Service
var wsLogger service.Logger
var svcRun bool
var cfgFile string
var hostAddr string

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of hashsrv:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "Options can be specified in the \"hashsrv.config\" file,\nor via environment variables prefixed with \"HASHSRV_\".\n")
	fmt.Fprintf(os.Stderr, "The location of the config file can be set with\nthe \"HASHSRV_CONFIG\" environment variable; otherwise\na standard set of locations is searched.\n")
}

func init() {
	const (
		name        = "HashSrv"
		displayName = "HashSrv"
		desc        = "HashSrv is a service for hashing, compressing, encrypting, and encoding things."
	)

	var help bool
	var ver bool
	var noisy bool
	var svcInstall bool
	var svcRemove bool
	var svcStart bool
	var svcStop bool

	flag.BoolVar(&noisy, "noisy", false, "Enable logging")
	flag.BoolVar(&help, "help", false, "Show command help")
	flag.BoolVar(&ver, "version", false, "Show version")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "Write CPU profile to file")
	flag.StringVar(&memprofile, "memprofile", "", "Write memory profile to file")
	flag.BoolVar(&svcInstall, "install", false, "Install HashSrv as a service")
	flag.BoolVar(&svcRemove, "remove", false, "Remove the HashSrv service")
	flag.BoolVar(&svcRun, "run", false, "Run HashSrv standalone (not as a service)")
	flag.BoolVar(&svcStart, "start", false, "Start the HashSrv service")
	flag.BoolVar(&svcStop, "stop", false, "Stop the HashSrv service")
	flag.StringVar(&hostAddr, "addr", ":9009", "Address to host the service on")
	flag.Parse()
	flagcfg.AddDefaults()
	flagcfg.Parse()
	flagenv.Prefix = "HASHSRV_"
	flagenv.Parse()

	if help {
		usage()
		os.Exit(0)
	}

	if ver {
		fmt.Printf("HashSrv Version %s\n", HASHSRV_VERSION)
		os.Exit(0)
	}

	var err error
	var i impl
	wsHashSrv, err = service.New(i, &service.Config{Name: name, DisplayName: displayName, Description: desc})
	if err != nil {
		log.Fatal(err)
	}
	wsLogger, err = wsHashSrv.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if svcInstall == true && svcRemove == true {
		log.Fatalln("Options -install and -remove cannot be used together.")
	} else if svcInstall == true {
		err = wsHashSrv.Install()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Service \"%s\" installed.\n", displayName)
		os.Exit(0)
	} else if svcRemove == true {
		err = wsHashSrv.Uninstall()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Service \"%s\" removed.\n", displayName)
		os.Exit(0)
	} else if svcStart == true {
		err = wsHashSrv.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Service \"%s\" started.\n", displayName)
		os.Exit(0)
	} else if svcStop == true {
		err = wsHashSrv.Stop()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Service \"%s\" stopped.\n", displayName)
		os.Exit(0)
	}

	if noisy == false {
		log.SetOutput(ioutil.Discard)
	}
}

type impl int

func (i impl) Start(s service.Service) error {
	// start
	go startWork()
	wsLogger.Info(fmt.Sprintf("Started HashSrv using config file \"%s\"", cfgFile))
	log.Printf("Started HashSrv using config file \"%s\"\n", cfgFile)
	return nil

}

func (i impl) Stop(s service.Service) error {
	// stop
	stopWork()
	wsLogger.Info("Stopped HashSrv")
	log.Println("Stopped HashSrv")
	return nil
}

func main() {
	var err error

	runtime.GOMAXPROCS(runtime.NumCPU())

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if svcRun == true {
		startWork()
		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, os.Interrupt, os.Kill)
		for {
			select {
			case event := <-sigChan:
				log.Print(event)
				switch event {
				case os.Interrupt, os.Kill: //SIGINT, SIGKILL
					return
				}
			}
		}
		stopWork()
	} else {
		err = wsHashSrv.Run()
		if err != nil {
			wsLogger.Error(err.Error())
			log.Println(err)
		}
	}
}

func startWork() {
	http.HandleFunc("/", root)
	go http.ListenAndServe(hostAddr, nil)
}

func stopWork() {
	// write memory profile if configured
	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Print(err)
		} else {
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}
}
