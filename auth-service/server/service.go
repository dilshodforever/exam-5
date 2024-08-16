package server

import (
	"fmt"
	"log"
	"net"

	"github.com/dilshodforever/5-oyimtixon/api"
	"github.com/dilshodforever/5-oyimtixon/api/handler"
	pb "github.com/dilshodforever/5-oyimtixon/genprotos/auth"
	pbu "github.com/dilshodforever/5-oyimtixon/genprotos/user"
	"github.com/dilshodforever/5-oyimtixon/service"
	"github.com/dilshodforever/5-oyimtixon/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Connect() (*gin.Engine, error) {
	db, err := postgres.NewPostgresStorage()
	if err != nil {
		log.Fatal("Error while connection on db: ", err.Error())
		return nil, err
	}
	liss, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal("Error while connection on tcp: ", err.Error())
		return nil, err
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, service.NewAuthService(db))
	pbu.RegisterUserServiceServer(s, service.NewUserService(db))
	log.Printf("server listening at %v", liss.Addr())
	go func() {
		if err := s.Serve(liss); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	UserConn, err := grpc.NewClient(fmt.Sprintf("auth-service%s", ":8085"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error while NEwclient: ", err.Error())
		return nil, err
	}
	defer UserConn.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	redisstorage := handler.NewInMemoryStorage(rdb)
	aus := service.NewAuthService(db)
	us := service.NewUserService(db)
	h := handler.NewHandler(aus, us, redisstorage)
	r := api.NewGin(h)
	return r, nil

}
