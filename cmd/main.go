package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/MosinEvgeny/authservice/internal/config"
	deliveryhttp "github.com/MosinEvgeny/authservice/internal/delivery/http"
	"github.com/MosinEvgeny/authservice/internal/repository/postgres"
	"github.com/MosinEvgeny/authservice/internal/service"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type EmailSender struct{}

func (sender *EmailSender) SendEmail(to, subject, body string) error {

	fmt.Println("MOCK Email sent to:", to) // TODO: Заменить на реальную отправку email
	fmt.Println("Subject:", subject)
	fmt.Println("Body:", body)
	return nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)

	}
	defer db.Close()

	emailSender := &EmailSender{}

	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)

	authService := service.NewAuthService(refreshTokenRepo, cfg, emailSender)

	handler := deliveryhttp.NewHandler(authService)

	router := mux.NewRouter()

	handler.InitRoutes(router)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: h2c.NewHandler(router, &http2.Server{}),
	}

	listener, err := net.Listen("tcp", cfg.ServerAddress)

	if err != nil {

		log.Fatal(err)

	}

	log.Printf("Server started on %s", cfg.ServerAddress)

	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {

		log.Fatalf("failed to start server %v", err)

	}

	log.Println("Shutting down server...")
}
