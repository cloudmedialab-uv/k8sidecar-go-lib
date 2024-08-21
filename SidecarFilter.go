package sidecar

// This file is responsible for defining the behavior of the sidecar filter, including forwarding requests to the next service in the chain and handling incoming requests.

// Importing the required packages.
import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// SidecarFilter is a struct that holds two potential functions: TriFunction and QuaFunction.
type SidecarFilter struct {
	TriFunction TriFunction
	QuaFunction QuaFunction
}

// FilterChain is a struct representing the chain of filters that a request goes through.
type FilterChain struct {
	Req *http.Request
	Res http.ResponseWriter
}

// Next method for the FilterChain advances to the next service by incrementing the port and sending the request to it.
func (chain *FilterChain) Next() {
	// Getting the port from the environment variable "PPORT" and incrementing it.
	port, _ := strconv.Atoi(os.Getenv("PPORT"))
	port++

	// Constructing the next service's URL.
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", "127.0.0.1", port),
		Path:   chain.Req.URL.Path,
	}

	// Creating a new request with the updated URL.
	req, err := http.NewRequest(chain.Req.Method, u.String(), chain.Req.Body)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header = chain.Req.Header

	// Sending the request and logging the response.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	for key, values := range res.Header {
		for _, value := range values {
			chain.Res.Header().Add(key, value)
		}
	}

	chain.Res.WriteHeader(res.StatusCode)

	body, err := io.Copy(chain.Res, res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Response copied successfully: %d bytes\n", body)
}

// ServeHTTP method implements the http.Handler interface. It processes the request, checks if it should invoke QuaFunction or TriFunction, and handles the request accordingly.
func (filter *SidecarFilter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Reading the original request body.
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return
	}
	req.Body.Close()

	// Creating a new reader for the original request body.
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Making a deep copy of the original request for the FilterChain.
	reqCopy := new(http.Request)
	*reqCopy = *req
	reqCopy.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Initializing the filter chain with the copied request and original response.
	chain := &FilterChain{Req: reqCopy, Res: res}

	value, exist := os.LookupEnv("PDISABLE")

	if exist {
		disable, err := strconv.ParseBool(value)

		if disable && err == nil {
			log.Println("SKIPPED sidecar")
			chain.Next()
			return
		}
	}

	// Checking if QuaFunction is defined, creating a cloud event from the request, and invoking QuaFunction.
	if filter.QuaFunction != nil {
		event, err := cloudevents.NewEventFromHTTPRequest(req)
		if err != nil {
			log.Println(err)
			return
		}

		filter.QuaFunction(req, res, *event, chain)
	} else {
		// If QuaFunction is not defined, invoking TriFunction.
		filter.TriFunction(req, res, chain)
	}
}

// Listen method starts the SidecarFilter server on the port specified by the "PPORT" environment variable.
func (filter *SidecarFilter) Listen() {
	port := os.Getenv("PPORT")
	http.Handle("/", filter)
	err := http.ListenAndServe("127.0.0.1:"+port, nil)
	if err != nil {
		log.Println(err)
	}
}
