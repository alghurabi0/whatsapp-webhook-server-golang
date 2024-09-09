package db

type Contact struct {
	WaId string `firestore:"-"`
	Name string `firestore:"name"`
}
