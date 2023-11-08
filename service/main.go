package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pritunl/pritunl-client-electron/service/auth"
	"github.com/pritunl/pritunl-client-electron/service/autoclean"
	"github.com/pritunl/pritunl-client-electron/service/config"
	"github.com/pritunl/pritunl-client-electron/service/constants"
	"github.com/pritunl/pritunl-client-electron/service/logger"
	"github.com/pritunl/pritunl-client-electron/service/profile"
	"github.com/pritunl/pritunl-client-electron/service/router"
	"github.com/pritunl/pritunl-client-electron/service/setup"
	"github.com/pritunl/pritunl-client-electron/service/tuntap"
	"github.com/pritunl/pritunl-client-electron/service/update"
	"github.com/pritunl/pritunl-client-electron/service/utils"
	"github.com/pritunl/pritunl-client-electron/service/watch"
	"github.com/pritunl/pritunl-client-electron/service/winsvc"
	"github.com/sirupsen/logrus"
)

func main() {
	install := flag.Bool("install", false, "run post install")
	uninstall := flag.Bool("uninstall", false, "run pre uninstall")
	clean := flag.Bool("clean", false, "clean up tuntap adapters")
	devPtr := flag.Bool("dev", false, "development mode")
	flag.Parse()

	if *install {
		setup.Install()
		return
	}

	if *uninstall {
		setup.Uninstall()
		return
	}

	if *clean {
		err := setup.TunTapClean(true)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	if *devPtr {
		constants.Development = true
	}

	err := config.Load()
	if err != nil {
		panic(err)
	}

	err = utils.PidInit()
	if err != nil {
		panic(err)
	}

	err = utils.InitTempDir()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("main: Failed to init temp dir")
		panic(err)
	}

	if runtime.GOOS == "darwin" {
		output, err := utils.ExecOutput("uname", "-r")
		if err == nil {
			macosVersion, err := strconv.Atoi(strings.Split(output, ".")[0])
			if err == nil && macosVersion < 20 {
				constants.Macos10 = true
			}
		}
	}

	logger.Init()

	logrus.WithFields(logrus.Fields{
		"version": constants.Version,
	}).Info("main: Service starting")

	go update.Check()

	defer func() {
		panc := recover()
		if panc != nil {
			logrus.WithFields(logrus.Fields{
				"stack": string(debug.Stack()),
				"panic": panc,
			}).Error("main: Panic")
			panic(panc)
		}
	}()

	err = auth.Init()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("main: Failed to init auth")
		panic(err)
	}

	err = autoclean.CheckAndClean()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("main: Failed to run check and clean")
		panic(err)
	}

	if runtime.GOOS == "windows" {
		if config.Config.DisableNetClean {
			logrus.Info("main: Network clean disabled")
		} else {
			err = tuntap.Clean()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("main: Failed to clear interfaces")
				err = nil
			}
		}
	}

	gin.SetMode(gin.ReleaseMode)

	watch.StartWatch()

	err = profile.Clean()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("main: Failed to clean profiles")
		panic(err)
	}

	routr := &router.Router{}
	routr.Init()

	go func() {
		retryCount := 0

		for {
			if constants.Interrupt {
				return
			}

			if retryCount > 2 {
				logrus.Error("main: Failed to recover server, closing")
				time.Sleep(5 * time.Second)
				panic(err)
			}

			err = routr.Run()
			if constants.Interrupt {
				return
			}
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("main: Server error")
			} else {
				logrus.Error("main: Unexpected server close")
			}

			if retryCount == 0 {
				go func() {
					time.Sleep(60 * time.Second)
					retryCount = 0
				}()
			}
			retryCount += 1

			newRoutr := &router.Router{}
			newRoutr.Init()
			routr = newRoutr

			time.Sleep(1 * time.Second)
		}
	}()

	profile.WatchSystemProfiles()

	if winsvc.IsWindowsService() {
		service := winsvc.New()

		err = service.Run()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("main: Service error")
			panic(err)
		}
	} else {
		sig := make(chan os.Signal, 100)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
	}

	constants.Interrupt = true

	routr.Shutdown()

	time.Sleep(250 * time.Millisecond)

	profile.Shutdown()

	prfls := profile.GetProfiles()
	for _, prfl := range prfls {
		prfl.StopBackground()
	}

	for _, prfl := range prfls {
		prfl.Wait()
	}

	if runtime.GOOS == "darwin" {
		_ = utils.ClearScutilConnKeys()
		_ = utils.RestoreScutilDns(true)
	}

	time.Sleep(750 * time.Millisecond)
}
