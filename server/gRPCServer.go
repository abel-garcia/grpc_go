package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/wackGarcia/protocol_buffers/book"
	"google.golang.org/grpc"
)

type server struct {
	syncProcess sync.Mutex
	book.UnimplementedAddressBookServer
	People map[string]*book.Person
}

// GRPCServer run server
func GRPCServer(port string) {
	//Run server
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gRPCServer := grpc.NewServer()

	book.RegisterAddressBookServer(gRPCServer, &server{})

	fmt.Printf("Server running on port %s \n", port)

	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// Methods called from client GRPC
func (s *server) Get(ctx context.Context, in *book.Person) (*book.Person, error) {
	s.syncProcess.Lock()
	person, ok := s.People[in.Email]
	s.syncProcess.Unlock()
	if !ok {
		return in, fmt.Errorf("Error: person not found")
	}

	return person, nil
}

func (s *server) Put(ctx context.Context, in *book.Person) (*book.Person, error) {
	people := make(map[string]*book.Person)
	people[in.GetEmail()] = in
	s.syncProcess.Lock()
	s.People = people
	s.syncProcess.Unlock()
	return in, nil
}

func (s *server) Del(ctx context.Context, in *book.Person) (*book.Person, error) {
	s.syncProcess.Lock()
	delete(s.People, in.Email)
	s.syncProcess.Unlock()
	return in, nil
}
