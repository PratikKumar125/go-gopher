package main

import (
	"context"
	"first/crons"
	"first/di"
	"first/models"
	"first/repositories"
	"first/routes"
	"first/tasks"
	"first/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

type MainPackage struct {
	cache       *utils.Cache
	pratikRepo  *repositories.UserRepository
  	asynqClient *utils.AsynqClient
  	asynqServer *utils.AsynqServer
	router 		*routes.Router
	cronRunner  *crons.CronRunnerStruct
}

func NewHandler(cache *utils.Cache, pratikRepo *repositories.UserRepository, asynqClient *utils.AsynqClient, asynqServer *utils.AsynqServer, router *routes.Router, cronRunner *crons.CronRunnerStruct) *MainPackage {
	return &MainPackage{cache: cache, pratikRepo: pratikRepo, asynqClient: asynqClient, asynqServer: asynqServer, router: router, cronRunner: cronRunner}
}

func (mp *MainPackage) InsertUser(c *gin.Context) {
	user := &models.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

    db_user_oid, err := mp.pratikRepo.CreateUser(context.Background(), user)
	if err != nil {
		fmt.Println("FAILED TO INSERT USER")
		return
	}

	task1, err := tasks.NewWelcomeEmailTask(db_user_oid)
	if err != nil {
		fmt.Println("unable to create email task")
		return
	}
	if _, err := mp.asynqClient.Client().Enqueue(
		task1,
		asynq.Queue("critical"),
	); err != nil {
		log.Fatal(err)
	}
	fmt.Println("INSERT USER", db_user_oid)
}

func main() {
	if err := di.InitDependencies(); err != nil {
		fmt.Println("Failed to initialize dependencies:", err)
		panic(err)
	}

	err := di.Container.Invoke(func(inj *di.Injected) {
	handler := NewHandler(inj.Utils.Cache, inj.Repositories.PratikRepo, inj.Utils.AsynqClientStruct, inj.Utils.AsynqServerStruct, inj.Router.Router, inj.Crons.CronRunner)

		//intializing the .env to os directly so that env vars can be accessed using os
		err := godotenv.Load(".env")
		if err != nil {
			panic("Failed to load env configuration")
		}
		app_port := os.Getenv("APP_PORT")
		fmt.Println("env value of key APP_PORT", app_port)

		// Start the Asynq server with the task handler
    	// If you want to have multiple workers for handling different types of tasks 
		//then you can fire two goroutines accordingly and similarly have created two 
		//different servers in the utils file

		// func() {
		// 	mux := asynq.NewServeMux()
		// 	mux.HandleFunc(tasks.TypeWelcomeEmail, inj.Utils.TaskHandlerStruct.HandleWelcomeEmailTask)
		// 	if err := inj.Utils.AsynqServerStruct.Server().Run(mux); err != nil {
		// 		log.Fatalf("Could not start Asynq server: %v", err)
		// 	}
		// }()

		//Initalize the console commands here
		// func () {
		// 	inj.Commands.DummyCommandStruct.RegisterDummyCommand()
		// } ()

		//Initialize all the CRON jobs here
		inj.Crons.CronRunner.RegisterCronJobs()

		//Initialize api router
		func() {
			handler.router.StartServer()
		}()
	})

	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
