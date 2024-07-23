package main

import (
	"first/crons"
	"first/di"
	"first/repositories/user_repository"
	"first/routes"
	"first/tasks"
	"first/utils"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type MainPackage struct {
	cache       *utils.Cache
	pratikRepo  *user_repository.UserRepository
  	asynqClient *utils.AsynqClient
  	asynqServer *utils.AsynqServer
	router 		*routes.Router
	cronRunner  *crons.CronRunnerStruct
	tasks *tasks.HandlerStruct
}

func NewHandler(cache *utils.Cache, pratikRepo *user_repository.UserRepository, asynqClient *utils.AsynqClient, asynqServer *utils.AsynqServer, router *routes.Router, cronRunner *crons.CronRunnerStruct, tasks *tasks.HandlerStruct) *MainPackage {
	return &MainPackage{cache: cache, pratikRepo: pratikRepo, asynqClient: asynqClient, asynqServer: asynqServer, router: router, cronRunner: cronRunner, tasks: tasks}
}

func main() {
	if err := di.InitDependencies(); err != nil {
		fmt.Println("Failed to initialize dependencies:", err)
		panic(err)
	}

	err := di.Container.Invoke(func(inj *di.Injected) {
		handler := NewHandler(inj.Utils.Cache, inj.Repositories.PratikRepo, inj.Utils.AsynqClientStruct, inj.Utils.AsynqServerStruct, inj.Router.Router, inj.Crons.CronRunner, inj.Tasks.Handler)

		//intializing the .env to os directly so that env vars can be accessed using os
		err := godotenv.Load(".env")
		if err != nil {
			panic("Failed to load env configuration")
		}
		app_port := os.Getenv("APP_PORT")
		fmt.Println(os.Getenv("JWT_SECRET"))
		fmt.Println("env value of key APP_PORT", app_port)

		// Start the Asynq server with the task handler
    	// If you want to have multiple workers for handling different types of tasks 
		//then you can fire two goroutines accordingly and similarly have created two 
		//different servers in the utils file

		// go func() {
		// 	mux := asynq.NewServeMux()
		// 	mux.HandleFunc(tasks.TypeWelcomeEmail, inj.Tasks.Handler.HandleWelcomeEmailTask)
		// 	if err := inj.Utils.AsynqServerStruct.Server().Run(mux); err != nil {
		// 		log.Fatalf("Could not start Asynq server: %v", err)
		// 	}
		// 	fmt.Println("Queue worker started")
		// }()

		//Initalize the console commands here
		// go func () {
		// 	inj.Commands.DummyCommandStruct.RegisterDummyCommand()
		// } ()

		//Initialize all the CRON jobs here
		func() {
			inj.Crons.CronRunner.RegisterCronJobs()
		}()

		//Initialize api router
		func() {
			handler.router.StartServer()
		}()
	})

	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
