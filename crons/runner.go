package crons

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

type CronRunnerStruct struct {
	TestCron TestCronStruct	
}

func NewCronRunner(testCron *TestCronStruct) *CronRunnerStruct {
	return &CronRunnerStruct{
		TestCron: *testCron,
	}
}

func (crs *CronRunnerStruct) RegisterCronJobs() {
	 c := cron.New(cron.WithSeconds())
	 c.AddFunc(TestCronTime, crs.TestCron.execute)
	 fmt.Println("Registered all CRON jobs")

	 c.Start()
	 select {}
}