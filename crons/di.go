package crons

import (
	"first/repositories/user_repository"

	"go.uber.org/dig"
)

type DependenciesHolder struct {
	dig.In
	TestCron *TestCronStruct
	CronRunner	*CronRunnerStruct
}

func RegisterDependencies(container *dig.Container) error {
	if err := container.Provide(func(UserRepo *user_repository.UserRepository) *TestCronStruct {
		return NewTestCronStruct(UserRepo)
	}); err != nil {
		return err
	}

	if err := container.Provide(func(TestCronStruct *TestCronStruct) *CronRunnerStruct {
		return NewCronRunner(TestCronStruct)
	}); err != nil {
		return err
	}
	return nil
}