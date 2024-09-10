package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/alghurabi0/whatsapp-webhook-server-golang/internal/models"
	"google.golang.org/api/option"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	contact  *models.ContactModel
	message  *models.MessageModel
}

func main() {
	addr := flag.String("addr", ":4002", "HTTP network address")
	credFile := flag.String("cred-file", "./internal/whatsapp-3a492-firebase-adminsdk-wd7lf-7f8138bbd2.json", "Path to the credentials file")
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := initDB(context.Background(), *credFile)
	if err != nil {
		errorLog.Fatalf("coldn't init db, err: %v\n", err)
	}

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		contact:  &models.ContactModel{DB: db},
		message:  &models.MessageModel{DB: db},
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
