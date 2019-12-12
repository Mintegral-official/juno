package builder

import (
	"log"
	"testing"
	"time"
)

func TestNewIndexBuilder(t *testing.T) {
	timer := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-timer.C:
			go func() {
				log.Println(time.Now())
			}()
		}
	}
}
