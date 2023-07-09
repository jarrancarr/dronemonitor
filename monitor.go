package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Location struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Altitude  float64 `json:"alt"`
	State     string  `json:"state"`
	Battery   float64 `json:"batt"`
}

var ctx = context.Background()

func main() {
	id := flag.Int("id", 9, "drone ID")
	redisHost := flag.String("redis", "localhost:6379", "Address of redis host")
	flag.Parse()

	redisClient := redis.NewClient(&redis.Options{
		Addr: *redisHost,
	})
	fmt.Printf("drone-%d\n", *id)
	subscriber := redisClient.Subscribe(ctx, fmt.Sprintf("drone-%d", *id))

	location := Location{}

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal([]byte(msg.Payload), &location); err != nil {
			panic(err)
		}

		fmt.Printf("Received message from " + msg.Channel + " channel.")
		fmt.Printf("%+v\n", location)
	}
}
