package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var queries = []string{
	"кроссовки асикс", "худи y2k", "айфон 17", "наушники маршал",
	"ремешок эпл вотч", "фильтр для воды", "духи том форд", "электрочайник",
	"сумка дизель", "кружка подарок папе", "мышка игровая", "набор ручек 100шт",
}

func main() {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	ctx := context.Background()
	js, _ := jetstream.New(nc)

	js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "SEARCH",
		Subjects: []string{"search.events"},
	})

	fmt.Println("Producing 1000 events...")
	for i := 0; i < 1000; i++ {
		event := map[string]any{
			"query":      queries[rand.Intn(len(queries))],
			"user_id":    fmt.Sprintf("u-%d", rand.Intn(100)),
			"session_id": fmt.Sprintf("s-%d", rand.Intn(500)),
			"timestamp":  time.Now().Unix(),
		}
		data, _ := json.Marshal(event)
		js.Publish(ctx, "search.events", data)
	}
	fmt.Println("Done! Now curl http://localhost:8080/api/v1/top?n=10")
}
