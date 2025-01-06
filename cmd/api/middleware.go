package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) rate(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var clients = make(map[string]*client)
	var mu sync.Mutex

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip_addr, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip_addr)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			ip_addr := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip_addr]; !found {
				clients[ip_addr] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			clients[ip_addr].lastSeen = time.Now()

			if !clients[ip_addr].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
