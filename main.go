package main

import (
	"SLON_tg_bot/src/app"
	"context"
	"log"
	"net/http"
	"os"
)

//func init() {
//	if err := godotenv.Load(); err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	token := os.Getenv("SLON_TOKEN")
	if token == "" {
		panic("empty token")
	}

	postgresConn := os.Getenv("PSQL_CONN")
	if postgresConn == "" {
		panic("empty psql connStr")
	}

	redisConn := os.Getenv("REDIS_CONN")
	if redisConn == "" {
		panic("empty redis connStr")
	}

	bot, err := app.NewBot(token, postgresConn, redisConn)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		bot.Bot.Start(ctx)
		log.Println("Bot started.")
	}()

	// in order to work with strange cloud.ru container system
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := ":8080"
	log.Printf("HTTP server is running on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
