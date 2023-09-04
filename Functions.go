package sidecar

// This file defines two custom function types for handling HTTP requests and responses in different scenarios.

// Importing the required packages.
import (
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// TriFunction is a type representing a function that takes in an HTTP request, an HTTP response writer, and a FilterChain.
type TriFunction func(req *http.Request, res http.ResponseWriter, chain *FilterChain)

// QuaFunction is a type representing a function that takes in an HTTP request, an HTTP response writer, a cloud event, and a FilterChain.
type QuaFunction func(*http.Request, http.ResponseWriter, cloudevents.Event, *FilterChain)
