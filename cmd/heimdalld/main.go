package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattmeyers/heimdall/http"
	"github.com/mattmeyers/heimdall/logger"
	"github.com/mattmeyers/heimdall/store/mem"
	"github.com/mattmeyers/heimdall/user"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := initializeFlags()

	logLevel, err := logger.ParseLevel(flags.logLevel)
	if err != nil {
		return err
	}

	logger, err := logger.NewLevelLogger(logLevel, nil)
	if err != nil {
		return err
	}

	userStore, err := mem.NewUserStore(mem.NewDB())
	if err != nil {
		return err
	}

	userService, err := user.NewService(userStore)
	if err != nil {
		return err
	}

	userController := &http.UserController{Service: *userService}

	s, err := http.NewServer(":8080", logger)
	if err != nil {
		return err
	}

	s.RegisterRoutes(userController)

	return s.ListenAndServe()
}

type flags struct {
	storeDriver string
	logLevel    string
}

func initializeFlags() flags {
	var fs flags

	flag.StringVar(&fs.storeDriver, "driver", "mem", "Database driver: mem")
	flag.StringVar(&fs.logLevel, "log-level", "info", "Min log level: debug, info, warn, error, fatal")

	flag.Parse()

	return fs
}
