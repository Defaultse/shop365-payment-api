package main

import (
	"context"
	"shop365-payment-api/internal/http"
)

func main() {
	srv := http.NewServer(context.Background(), ":8080")

	err := srv.Run()
	if err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}
