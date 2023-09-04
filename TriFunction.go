package sidecar

import "net/http"

type TriFunction func(req *http.Request, res http.ResponseWriter, chain *FilterChain)
