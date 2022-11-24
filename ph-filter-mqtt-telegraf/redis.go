package main

import (
	"context"
	"fmt"
	"log"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	goredis "github.com/go-redis/redis/v8"
)

type redisClient struct {
	client *goredis.Client
}

func createRedisClientConnection(addr string, pass string, db int) *redisClient {
	//handler := rejson.NewReJSONHandler()
	//flag.Parse()

	// GoRedis Client
	cli := goredis.NewClient(&goredis.Options{Addr: addr, Password: pass, DB: db})

	// Se scommentata restituisce "client is closed"
	// defer func() {
	// 	if err := cli.FlushAll(context.Background()).Err(); err != nil {
	// 		log.Fatalf("goredis - failed to flush: %v", err)
	// 	}
	// 	if err := cli.Close(); err != nil {
	// 		log.Fatalf("goredis - failed to communicate to redis-server: %v", err)
	// 	}
	// }()
	//handler.SetGoRedisClient(cli)
	redisC := redisClient{client: cli}
	return &redisC
}

func (redisC *redisClient) saveToRedis(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, result interface{}) {
	err := redisC.client.Set(context.Background(), "PhTemperature", data.(map[string]interface{})["Ph"].(PhData).Temperature, 0).Err()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Failed to Set")
		return
	}
	err = redisC.client.MSet(context.Background(), "Ph", data.(map[string]interface{})["Ph"].(PhData).PhValue, 0).Err()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Failed to Set")
		return
	}
	return true, nil
}
