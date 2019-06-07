package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	_ "github.com/lib/pq"
)

const Port = "9999"
const BackendEndpoint = "54.199.144.215:8888/backend"

type frontHandler struct{}

func (h *frontHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("front requested, reponding with HTTP 200")

	client := xray.Client(&http.Client{})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", BackendEndpoint), nil)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		fmt.Fprint(writer, "Error New Request")
		return
	}

	resp, err := client.Do(req.WithContext(request.Context()))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		fmt.Fprint(writer, "Error Request")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		fmt.Fprint(writer, "Error Read body")
		return
	}

	writer.WriteHeader(resp.StatusCode)
	fmt.Fprint(writer, string(body))
}

func main() {
	fmt.Println("starting server listening on port " + Port)
	xraySegmentNamer := xray.NewFixedSegmentNamer("front-colorteller")
	http.Handle("/front", xray.Handler(xraySegmentNamer, &frontHandler{}))
	http.ListenAndServe(":"+Port, nil)
}
