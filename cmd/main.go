package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"nexussiege/internal/server"
)

func main() {
	// Criar servidor
	srv := server.NewServer()

	// Canal para sinais de shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	errChan := make(chan error, 1)
	go func() {
		log.Println("Iniciando Nexus Siege...")
		if err := srv.Start(":8080"); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Aguardar sinal ou erro
	select {
	case err := <-errChan:
		log.Fatalf("Erro no servidor: %v", err)
	case <-sigChan:
		log.Println("Shutdown iniciado...")
	}

	// Shutdown gracioso
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Erro no shutdown: %v", err)
	}

	log.Println("Servidor parado.")
}
