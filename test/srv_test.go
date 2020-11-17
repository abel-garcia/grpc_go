package test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/wackGarcia/grpc_go/server"
	"github.com/wackGarcia/protocol_buffers/book"
	"google.golang.org/grpc"
)

var (
	setupOnce sync.Once
	client    book.AddressBookClient
	person    book.Person
)

func setup() {
	//Read config.json and set varibles
	viper.SetConfigFile("../config.json")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error in config file: %s", err))
	}

	target := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
	port := fmt.Sprintf(":%s", viper.GetString("server.port"))

	// running server test
	go server.GRPCServer(port)

	// running client test
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("grpc.Dial: %v", err)
	}

	client = book.NewAddressBookClient(conn)

	person = book.Person{
		Id:       1,
		Name:     "Abel",
		LastName: "Garcia",
		Age:      28,
		Phones: []*book.Person_PhoneNumber{
			{Number: "555-4321", Type: book.Person_HOME},
		},
		Email: "abelgarcia38348@gmail.com",
	}

	time.Sleep(3 * time.Second)
}

func TestGet(t *testing.T) {

	setupOnce.Do(setup)

	if _, err := client.Put(context.Background(), &person); err != nil {
		t.Fatalf("Error in put: %v", err)
	}

	_, err := client.Get(context.Background(), &book.Person{Email: "abelgarcia38348@gmail.com"})

	if err != nil {
		t.Fatalf("Error in get: %v", err)
	}
}

func BenchmarkDataRaces(b *testing.B) {
	setupOnce.Do(setup)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ne1, _ := client.Get(context.Background(), &book.Person{Email: "abelgarcia38348@gmail.com"})
			ne2, _ := client.Put(context.Background(), &person)
			ne3, _ := client.Get(context.Background(), &book.Person{Email: "abelgarcia38348@gmail.com"})
			ne4, _ := client.Del(context.Background(), &book.Person{Email: "abelgarcia38348@gmail.com"})
			ne5, _ := client.Get(context.Background(), &book.Person{Email: "abelgarcia38348@gmail.com"})

			_ = ne1
			_ = ne2
			_ = ne3
			_ = ne4
			_ = ne5
		}

	})
}

func BenchmarkReadOnly(b *testing.B) {
	setupOnce.Do(setup)
	client.Put(context.Background(), &person)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ne1, _ := client.Get(context.Background(), &book.Person{Email: "abelgarcia38348@gmail.com"})
			_ = ne1
		}

	})
}

func BenchmarkWriteOnly(b *testing.B) {
	setupOnce.Do(setup)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ne1, _ := client.Put(context.Background(), &person)
			ne1, _ = client.Del(context.Background(), &person)
			_ = ne1
		}

	})
}
