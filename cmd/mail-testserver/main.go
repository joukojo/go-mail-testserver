package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joukojo/go-mail-testserver/internal/commonssmtp"
	httpapi "github.com/joukojo/go-mail-testserver/internal/httpapi"
)

func main() {
	fmt.Println("Mailserver")

	storage := httpapi.NewStorage()

	smtpAddr := getenv("SMTP_ADDR", ":1025")
	httpAddr := getenv("HTTP_ADDR", ":8025")

	fmt.Printf("Starting SMTP server at %s\n", smtpAddr)
	smtpServer := commonssmtp.NewSmtpServer(storage, smtpAddr)
	fmt.Printf("Starting HTTP server at %s\n", httpAddr)
	apiServer := httpapi.New(httpAddr, storage)

	go func() {
		if err := smtpServer.Start(); err != nil {
			fmt.Printf("SMTP server error: %v\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		if err := apiServer.Start(); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Printf("Starting HTTP server at %s\n", httpAddr)

	// Graceful shutdown on SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
