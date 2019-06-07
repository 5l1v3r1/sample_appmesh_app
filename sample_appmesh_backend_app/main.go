package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	_ "github.com/lib/pq"
)

const Port = "8888"

type backendHandler struct{}

func (h *backendHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("backend requested, reponding with HTTP 200")
	db, err := xray.SQL("postgres", "user=dev password=secret dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Printf("%+v", err)
	}
	defer db.Close()

	var res int
	ctx, _ := xray.BeginSegment(context.Background(), "test")
	err = db.QueryRow(ctx, "SELECT 200").Scan(&res)
	if err != nil {
		fmt.Printf("failed to error %+v\n", err)
	}
	writer.WriteHeader(res)
	fmt.Fprint(writer, "Hello, World!\n")
}

func main() {
	fmt.Println("starting server listening on port " + Port)
	xraySegmentNamer := xray.NewFixedSegmentNamer("backend-colorteller")
	http.Handle("/backend", xray.Handler(xraySegmentNamer, &backendHandler{}))
	http.ListenAndServe(":"+Port, nil)
}
