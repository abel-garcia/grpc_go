package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/wackGarcia/grpc_go/client"
	"github.com/wackGarcia/grpc_go/server"
)

func main() {

	//Read config.json and set varibles
	viper.SetConfigFile("config.json")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error in config file: %s", err))
	}

	target := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))
	port := fmt.Sprintf(":%s", viper.GetString("server.port"))

	go func() {
		fmt.Println("waitting for server...")
		time.Sleep(5 * time.Second)
		client.GRPClient(target)
	}()

	server.GRPCServer(port)
}
