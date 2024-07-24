package crons

import (
	"first/repositories/user_repository"
	"fmt"
)

type TestCronStruct struct {
	UserRepo *user_repository.UserRepository
}

func NewTestCronStruct(repo *user_repository.UserRepository) *TestCronStruct{
	fmt.Println("Test Cron Job Intialized")
	return &TestCronStruct{
		UserRepo: repo,
	}
}

var (TestCronTime = "*/30 * * * * *")

func (tcs *TestCronStruct) execute() {
    // res, err := tcs.UserRepo.FindAll(context.Background())
    // if err != nil {
    //     fmt.Printf("Error fetching users: %v\n", err)
    //     return
    // }
    // fmt.Printf("Users: %+v\n", res)
    fmt.Println("TEST CRON EXECUTED SUCCESSFULLY")
}
