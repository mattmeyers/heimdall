package main

import (
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
	logger, err := logger.NewLevelLogger(logger.LevelDebug, nil)
	if err != nil {
		return err
	}

	userStore, err := mem.NewUserStore()
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
