package main

import (
	"context"
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
	router *routes.Router 
}

func NewHandler(cache *utils.Cache, pratikRepo *repositories.UserRepository, asynqClient *utils.AsynqClient, asynqServer *utils.AsynqServer, router *routes.Router) *MainPackage {
	return &MainPackage{cache: cache, pratikRepo: pratikRepo, asynqClient: asynqClient, asynqServer: asynqServer, router: router}
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
		return
	}

	err := di.Container.Invoke(func(inj *di.Injected) {
	handler := NewHandler(inj.Utils.Cache, inj.Repositories.PratikRepo, inj.Utils.AsynqClientStruct, inj.Utils.AsynqServerStruct, inj.Router.Router)

		//Initalize the console commamds here
		// func () {
		// 	inj.Commands.DummyCommandStruct.RegisterDummyCommand()
		// } ()

		//intializing the .env to os directly so that env vars can be accessed using os
		err := godotenv.Load(".env")
		if err != nil {
			panic("Failed to load env configuration")
		}
		//Read from the config now directly
		helo_world := os.Getenv("HELLO_WORLD")
		fmt.Println("env hello_world", helo_world)

		//Initialize api router
		func() {
			handler.router.StartServer()
		}()

		// Start the Asynq server with the task handler
    // If you want to have multiple workers for handling different types of tasks then you can fire two goroutines accordingly and similarly have created two different servers in the utils file
		func() {
			mux := asynq.NewServeMux()
			mux.HandleFunc(tasks.TypeWelcomeEmail, inj.Utils.TaskHandlerStruct.HandleWelcomeEmailTask)
			if err := inj.Utils.AsynqServerStruct.Server().Run(mux); err != nil {
				log.Fatalf("Could not start Asynq server: %v", err)
			}
		}()

		//Initialize all the CRON jobs here
		// inj.Crons.CrocnRunner.RegisterCronJobs();
	})

	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
