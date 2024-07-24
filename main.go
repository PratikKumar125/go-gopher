package main

import (
	"first/crons"
	"first/di"
	"first/repositories/user_repository"
	"first/routes"
	"first/tasks"
	"first/utils"
	"fmt"
)

type MainPackage struct {
	cache       *utils.Cache
	userRepo  *user_repository.UserRepository
  	asynqClient *utils.AsynqClient
  	asynqServer *utils.AsynqServer
	router 		*routes.Router
	cronRunner  *crons.CronRunnerStruct
	tasks *tasks.HandlerStruct
}

func NewHandler(cache *utils.Cache, userRepo *user_repository.UserRepository, asynqClient *utils.AsynqClient, asynqServer *utils.AsynqServer, router *routes.Router, cronRunner *crons.CronRunnerStruct, tasks *tasks.HandlerStruct) *MainPackage {
	return &MainPackage{cache: cache, userRepo: userRepo, asynqClient: asynqClient, asynqServer: asynqServer, router: router, cronRunner: cronRunner, tasks: tasks}
}

func main() {
	//the .env file will be loaded in the InitDependencies() function itself
	if err := di.InitDependencies(); err != nil {
		fmt.Println("Failed to initialize dependencies:", err)
		panic(err)
	}

	err := di.Container.Invoke(func(inj *di.Injected) {
		handler := NewHandler(inj.Utils.Cache, inj.Repositories.UserRepo, inj.Utils.AsynqClientStruct, inj.Utils.AsynqServerStruct, inj.Router.Router, inj.Crons.CronRunner, inj.Tasks.Handler)

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
