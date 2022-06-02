package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattmeyers/heimdall/auth"
	"github.com/mattmeyers/heimdall/client"
	"github.com/mattmeyers/heimdall/http"
	"github.com/mattmeyers/heimdall/store"
	"github.com/mattmeyers/heimdall/store/sqlite"
	"github.com/mattmeyers/heimdall/user"
	"github.com/mattmeyers/level"
	_ "modernc.org/sqlite"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := initializeFlags()

	logLevel, err := level.ParseLevel(flags.logLevel)
	if err != nil {
		return err
	}

	logger, err := level.NewBasicLogger(logLevel, nil)
	if err != nil {
		return err
	}

	var ss stores
	switch flags.storeDriver {
	case "mem":
		ss, err = getSqliteStores("file::memory:", false)
	case "sqlite":
		ss, err = getSqliteStores("file:db/data/heimdall-dev.db?mode=rwc", flags.noMigrate)
	default:
		return errors.New("unknown driver")
	}

	if err != nil {
		return err
	}

	logger.Info("Using DB driver: %s", flags.storeDriver)

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

	authService, err := auth.NewService(
		ss.userStore,
		ss.clientStore,
		ss.authCodeStore,
		auth.JWTSettings{
			Issuer:     "heimdall",
			Lifespan:   3600,
			SigningKey: "so-secret-wow",
			Algorithm:  auth.HMAC256Algorithm,
		},
	)
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
	noMigrate   bool
}

func initializeFlags() flags {
	var fs flags

	flag.StringVar(&fs.storeDriver, "driver", "mem", "Database driver: mem, sqlite")
	flag.BoolVar(&fs.noMigrate, "no-migrate", false, "Prevent migrating db. Ignored for mem driver.")
	flag.StringVar(&fs.logLevel, "log-level", "info", "Min log level: debug, info, warn, error, fatal")

	flag.Parse()

	return fs
}

type stores struct {
	userStore     store.UserStore
	clientStore   store.ClientStore
	authCodeStore store.AuthCodeStore
}

func getSqliteStores(dsn string, noMigrate bool) (stores, error) {
	db, err := sqlite.NewDB(dsn)
	if err != nil {
		return stores{}, err
	}

	if !noMigrate {
		driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			return stores{}, err
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://./db/migrations/sqlite",
			"sqlite", driver)
		if err != nil {
			return stores{}, err
		}

		err = m.Up()
		if err != migrate.ErrNoChange && err != nil {
			return stores{}, err
		}
	}

	userStore, err := sqlite.NewUserStore(db)
	if err != nil {
		return stores{}, err
	}

	clientStore, err := sqlite.NewClientStore(db)
	if err != nil {
		return stores{}, err
	}

	authCodeStore, err := sqlite.NewAuthCodeStore(db)
	if err != nil {
		return stores{}, err
	}

	return stores{userStore: userStore, clientStore: clientStore, authCodeStore: authCodeStore}, nil
}
