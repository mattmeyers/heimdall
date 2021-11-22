package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/mattmeyers/heimdall/auth"
	"github.com/mattmeyers/heimdall/client"
	"github.com/mattmeyers/heimdall/http"
	"github.com/mattmeyers/heimdall/logger"
	"github.com/mattmeyers/heimdall/store"
	"github.com/mattmeyers/heimdall/store/mem"
	"github.com/mattmeyers/heimdall/store/sqlite"
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

	var ss stores
	switch flags.storeDriver {
	case "mem":
		ss, err = getMemStores()
	case "sqlite":
		ss, err = getSqliteStores("file:db/data/heimdall-dev.db?mode=rw")
	default:
		return errors.New("unknown driver")
	}

	if err != nil {
		return err
	}

	userService, err := user.NewService(ss.userStore)
	if err != nil {
		return err
	}

	userController := &http.UserController{Service: *userService}

	clientService, err := client.NewService(ss.clientStore)
	if err != nil {
		return err
	}

	clientController := &http.ClientController{Service: *clientService}

	authService, err := auth.NewService(ss.userStore)
	if err != nil {
		return err
	}

	authController := &http.AuthController{Service: *authService}

	s, err := http.NewServer(":8080", logger)
	if err != nil {
		return err
	}

	s.RegisterRoutes(userController, clientController, authController)

	return s.ListenAndServe()
}

type flags struct {
	storeDriver string
	logLevel    string
}

func initializeFlags() flags {
	var fs flags

	flag.StringVar(&fs.storeDriver, "driver", "mem", "Database driver: mem, sqlite")
	flag.StringVar(&fs.logLevel, "log-level", "info", "Min log level: debug, info, warn, error, fatal")

	flag.Parse()

	return fs
}

type stores struct {
	userStore   store.UserStore
	clientStore store.ClientStore
}

func getMemStores() (stores, error) {
	db := mem.NewDB()

	userStore, err := mem.NewUserStore(db)
	if err != nil {
		return stores{}, err
	}

	clientStore, err := mem.NewClientStore(db)
	if err != nil {
		return stores{}, err
	}

	return stores{userStore: userStore, clientStore: clientStore}, nil
}

func getSqliteStores(dsn string) (stores, error) {
	db, err := sqlite.NewDB(dsn)
	if err != nil {
		return stores{}, err
	}

	userStore, err := sqlite.NewUserStore(db)
	if err != nil {
		return stores{}, err
	}

	clientStore, err := sqlite.NewClientStore(db)
	if err != nil {
		return stores{}, err
	}

	return stores{userStore: userStore, clientStore: clientStore}, nil
}
