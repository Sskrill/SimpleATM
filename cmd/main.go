package main

import (
	"fmt"
	"github.com/Sskrill/SimpleATM/internal/service"
	"github.com/Sskrill/SimpleATM/internal/transport"
	"log"
	"net/http"
)

func main() {
	db := service.NewData()
	srvc := service.NewService(db)
	h := transport.NewHandler(srvc)
	server := &http.Server{Addr: fmt.Sprintf(":8080"), Handler: h.CreateRouter()}

	fmt.Println("Server started ")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
