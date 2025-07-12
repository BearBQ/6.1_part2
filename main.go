package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func main() {
	ctx := context.Background()
	url := []string{"https://catfact.ninja/fact", "https://ya.ru", "https://google.com"}
	i := 1
	var wg sync.WaitGroup
	for _, val := range url {
		wg.Add(1)
		id := "abc-" + strconv.Itoa(i)
		go fetchWithContext(ctx, val, id, &wg)
		i++
	}
	wg.Wait()

}

func fetchWithContext(ctx context.Context, url string, id string, wg *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	defer wg.Done()
	ctx = context.WithValue(ctx, requestIDKey, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("ошибка при создании запроса")
	}
	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Запрос отменён: %v, requestID=%s", ctx.Err(), id)
		} else {
			log.Printf("Ошибка при выполнении запроса: %v, requestID=%s\n", err, ctx.Value(requestIDKey).(string))
		}
		return
	}

	// Выводим результат
	fmt.Printf("requestID^ %s Статус: %s\n", ctx.Value(requestIDKey).(string), resp.Status)

}
