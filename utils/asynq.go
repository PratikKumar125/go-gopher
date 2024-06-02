package utils

import (
	"fmt"

	"github.com/hibiken/asynq"
)

type AsynqClient struct {
	client *asynq.Client
}

func NewAsynqClient() *AsynqClient {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: "localhost:6379", // redis server address
	})
	fmt.Println("ASYNQ CLIENT INITIALIZED")
	return &AsynqClient{client: client}
}

func (c *AsynqClient) Client() *asynq.Client {
	return c.client
}

type AsynqServer struct {
	server *asynq.Server
}

func NewAsynqServer() *AsynqServer {
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: "localhost:6379",
		},
		asynq.Config{
			Concurrency: 10, // number of concurrent workers
			Queues: map[string]int{
				"critical": 6, // processed 60% of the time
				"default":  3, // processed 30% of the time
				"low":      1, // processed 10% of the time
			},
		},
	)
	fmt.Println("ASYNQ SERVER INITIALIZED")
	return &AsynqServer{server: server}
}

func (s *AsynqServer) Server() *asynq.Server {
	return s.server
}
