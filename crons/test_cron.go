package crons

import (
	"first/repositories"
	"fmt"
)

type TestCronStruct struct {
	UserRepo *repositories.UserRepository
}

func NewTestCronStruct(repo *repositories.UserRepository) *TestCronStruct{
	fmt.Println("Test Cron Job Intialized")
	return &TestCronStruct{
		UserRepo: repo,
	}
}

var (TestCronTime = "*/30 * * * * *")

func (tcs *TestCronStruct) execute() {
	// newUser := &models.User{
	// 	Name: "CRON USER",
	// 	Email: "cron@cron.com",
	// }
	// tcs.UserRepo.CreateUser(nil, newUser)
	fmt.Println("TEST CRON EXECUTED SUCCESSFULLY")
}