package sidecar

import (
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type TriFunction func(req *http.Request, res http.ResponseWriter, chain *FilterChain)

type QuaFunction func(*http.Request, http.ResponseWriter, cloudevents.Event, *FilterChain)
