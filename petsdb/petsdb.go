package petsdb

import (

	"context"
	"fmt"
	"log"
	"os"
	"time"
	"cloud.google.com/go/datastore"
)

var projectID string

// Pet model stored in Datastore
type Pet struct {
	Added   time.Time `datastore:"added"`
	Caption string    `datastore:"caption"`
	Email   string    `datastore:"email"`
	Image   string    `datastore:"image"`
	Likes   int       `datastore:"likes"`
	Owner   string    `datastore:"owner"`
	Petname string    `datastore:"petname"`
	Id      string     // The ID used in the datastore.
}

func getDataStoreClient()(context.Context, *datastore.Client) {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}
	return ctx, client
}

// GetPets Returns all pets from datastore ordered by likes in Desc Order
func GetPets() ([]Pet, error) {
	var pets []Pet
	ctx, client := getDataStoreClient()

	// Create a query to fetch all Pet entities".
	query := datastore.NewQuery("Pet") //.Order("-likes")
	keys, err := client.GetAll(ctx, query, &pets)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Set the id field on each Task from the corresponding key.
	for i, key := range keys {
		pets[i].Id = key.Name
	}

	client.Close()
	return pets, nil
}

func GetPetById(petId string)(*Pet, error){
	pet := &Pet{}
	ctx, client := getDataStoreClient()
	key := datastore.NameKey("Pet", petId, nil)
	err := client.Get(ctx, key, pet)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	client.Close()
	return pet, nil
}

func CreatePet(pet Pet){
	ctx, client := getDataStoreClient()
	key := datastore.IncompleteKey("Pet", nil)
	_, err := client.Put(ctx, key, &pet)
	if err != nil {
		fmt.Println(err)
	}
	client.Close()
}
