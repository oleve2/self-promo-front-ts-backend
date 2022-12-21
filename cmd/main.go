package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"self_promo_back/cmd/app"

	dbaseServ "self_promo_back/pkg/dbase"
	emailServ "self_promo_back/pkg/email"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = "9999"
	sqlitePath  = "./selfpromo.db"
)

func main() {
	// read env variables
	err := godotenv.Load()
	if err != nil {
		//log.Fatal("Error loading .env file")
		log.Println("Error loading .env file")
	}

	// main app settings
	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	// email service env variables
	emailFrom, ok := os.LookupEnv("APP_email_from")
	if !ok {
		log.Fatal("email_from not set")
	}
	emailTo, ok := os.LookupEnv("APP_email_to")
	if !ok {
		log.Fatal("email_to not set")
	}
	smtpHost, ok := os.LookupEnv("APP_smtpHost")
	if !ok {
		log.Fatal("smtpHost not set")
	}
	smtpPort, ok := os.LookupEnv("APP_smtpPort")
	if !ok {
		log.Fatal("smtpPort not set")
	}
	password, ok := os.LookupEnv("APP_password")
	if !ok {
		log.Fatal("password not set")
	}

	// access token
	accesstoken, ok := os.LookupEnv("APP_accessToken")
	if !ok {
		log.Fatal("access token not set")
	}

	// start server
	if err := execute(
		net.JoinHostPort(host, port),
		emailFrom, emailTo, smtpHost, smtpPort, password, accesstoken); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(
	addr string,
	emailFrom string, emailTo string, smtpHost string, smtpPort string, password string,
	accesstoken string,
) error {
	// роутер
	mux := chi.NewRouter()

	// email service
	EmailSvc := emailServ.NewService(emailFrom, emailTo, smtpHost, smtpPort, password)

	// database service
	dbSqlite, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbaseSvc := dbaseServ.NewService(dbSqlite)
	dbaseSvc.Init()

	// backend http server
	application := app.NewServer(
		mux,
		EmailSvc,
		dbaseSvc,
		accesstoken,
	)

	// init app
	err = application.Init()
	if err != nil {
		log.Println(err)
		return err
	}
	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}

	fmt.Printf("server started on http://%s\n", addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
