package crons

import (
	"context"
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
    res, err := tcs.UserRepo.FindAll(context.Background())
    if err != nil {
        fmt.Printf("Error fetching users: %v\n", err)
        return
    }
    fmt.Printf("Users: %+v\n", res)
    fmt.Println("TEST CRON EXECUTED SUCCESSFULLY")
}
