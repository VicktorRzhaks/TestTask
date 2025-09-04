package main
import (
	"context"
	"fmt"
	"sync"
	"time"
)
func Request(id int) {
	fmt.Printf("Request %02d started  at %s\n", id, time.Now().Format("15:04:05.000"))
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Request %02d finished at %s\n", id, time.Now().Format("15:04:05.000"))
}
func main() {
	const (
		N      = 10
		M      = 3
		period = 200 * time.Millisecond
	)
	jobs := make(chan int)
	limiter := make(chan struct{}, 1) // токен-бакет 
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Генератор токенов: не чаще 1/period
	go func() {
		t := time.NewTicker(period)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return 
			case <-t.C:
				select {
				case limiter <- struct{}{}:
				default:
				}
			}
		}
	}()
	var wg sync.WaitGroup
	// Пул воркеров
	for i := 0; i < M; i++ {
		go func() {
			for id := range jobs {
				func(id int) {
					defer wg.Done() // гарантируем баланс на каждый взятый job
					// ждeм токен или отмену
					select {
					case <-limiter:
						Request(id)
					case <-ctx.Done():
						// отмена — выходим без request
						return
					}
				}(id)
			}
		}()
	}
	// аdd перед отправкой, закрываем jobs в том же потоке
	for i := 0; i < N; i++ {
		wg.Add(1) // учет задачи до отправки
		jobs <- i
	}
	close(jobs)
	// ждем завершения всех задач и останавливаем генератор токенов
	wg.Wait()
	cancel()
}
