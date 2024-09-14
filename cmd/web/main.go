package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
	"google.golang.org/api/option"
)

type application struct {
	infoLog         *log.Logger
	errorLog        *log.Logger
	templateCache   map[string]*template.Template
	contact         *models.ContactModel
	message         *models.MessageModel
	token           string
	phone_number_id string
}

func main() {
	addr := flag.String("addr", ":4002", "HTTP network address")
	credFile := flag.String("cred-file", "./internal/whatsapp-3a492-firebase-adminsdk-wd7lf-7f8138bbd2.json", "Path to the credentials file")
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	db, err := initDB(context.Background(), *credFile)
	if err != nil {
		errorLog.Fatalf("coldn't init db, err: %v\n", err)
	}
	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		errorLog.Fatal("empty access token")
	}
	phone_number_id := os.Getenv("PHONE_NUMBER_ID")
	if phone_number_id == "" {
		errorLog.Fatal("empty phone number id")
	}

	app := &application{
		infoLog:         infoLog,
		errorLog:        errorLog,
		templateCache:   templateCache,
		contact:         &models.ContactModel{DB: db},
		message:         &models.MessageModel{DB: db},
		token:           token,
		phone_number_id: phone_number_id,
	}
	srv := http.Server{
		Handler:  app.routes(),
		ErrorLog: errorLog,
		Addr:     *addr,
	}
	infoLog.Printf("app running on 4002\n")
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func initDB(ctx context.Context, credFile string) (*firestore.Client, error) {
	opt := option.WithCredentialsFile(credFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	//TODO - ping the database to check if it's connected
	docRef := firestoreClient.Collection("ping").Doc("test")
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := docSnapshot.DataTo(&data); err != nil {
		return nil, err
	}
	expectedValue := "pong"
	if value, ok := data["ping"].(string); !ok || value != expectedValue {
		return nil, fmt.Errorf("ping test failed, expected %s, got %s", expectedValue, value)
	}

	return firestoreClient, nil
}
