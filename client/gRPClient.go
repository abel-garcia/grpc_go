package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wackGarcia/protocol_buffers/book"
	"google.golang.org/grpc"
)

// GRPClient run client grpc example
func GRPClient(target string) {

	// Run client
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		panic(fmt.Errorf("did not connnect: %v", err))
	}

	defer conn.Close()

	client := book.NewAddressBookClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	// Example put person
	fmt.Println("Saving person..")
	putPerson(ctx, client)

	// Example get person
	fmt.Println("Reading person..")
	printPerson(ctx, client, "abelgarcia38348@gmail.com")

}

//Call methods from server
func putPerson(ctx context.Context, client book.AddressBookClient) {
	// Call method Put from serve
	r, err := client.Put(ctx, &book.Person{
		Id:       1,
		Name:     "Abel",
		LastName: "Garcia",
		Age:      28,
		Phones: []*book.Person_PhoneNumber{
			{Number: "555-4321", Type: book.Person_HOME},
		},
		Email: "abelgarcia38348@gmail.com",
	})

	if err != nil {
		log.Fatalf("could not put: %v", err)
	}

	fmt.Printf("Person %s, was saved \n", r.GetEmail())
}

func printPerson(ctx context.Context, client book.AddressBookClient, email string) {
	// Call method Get from server
	res, err := client.Get(ctx, &book.Person{
		Email: email,
	})

	if err != nil {
		log.Fatalf("could not get: %v", err)
	}

	fmt.Printf("ID: %d \n", res.GetId())
	fmt.Printf("Emial: %s \n", res.GetEmail())
	fmt.Printf("Name: %s \n", res.GetName())
	fmt.Printf("Last Name: %s ", res.GetLastName())
}
