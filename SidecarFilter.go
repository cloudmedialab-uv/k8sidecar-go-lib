package sidecar

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

type SidecarFilter struct {
	TriFunction TriFunction
	QuaFunction QuaFunction
}

type FilterChain struct {
	req *http.Request
	res http.ResponseWriter
}

func (chain *FilterChain) Next() {
	port, _ := strconv.Atoi(os.Getenv("PPORT"))
	port++
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", "127.0.0.1", port),
		Path:   chain.req.URL.Path,
	}

	req, err := http.NewRequest(chain.req.Method, u.String(), chain.req.Body)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header = chain.req.Header

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(body))

}

func (filter *SidecarFilter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Leemos el cuerpo original en un slice de bytes
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return
	}
	req.Body.Close() // Asegúrate de cerrar el body original cuando hayas terminado con él

	// Creamos un nuevo reader para el cuerpo de la solicitud, que se puede leer en filter.Function
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Creamos una copia de la solicitud para usar en FilterChain
	reqCopy := new(http.Request)
	*reqCopy = *req // copiamos todos los campos
	// Creamos un nuevo reader para el cuerpo de la solicitud en la copia, que se puede leer en chain.Next()
	reqCopy.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	chain := &FilterChain{req: reqCopy, res: res}

	value, exist := os.LookupEnv("PDISABLE")

	if exist {
		disable, err := strconv.ParseBool(value)

		if disable && err == nil {
			chain.Next()
		}
	}

	if filter.QuaFunction != nil {
		event, err := cloudevents.NewEventFromHTTPRequest(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		filter.QuaFunction(req, res, *event, chain)
	} else {

		filter.TriFunction(req, res, chain)
	}
}

func (filter *SidecarFilter) Listen() {
	port := os.Getenv("PPORT")
	http.Handle("/", filter)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err)
	}
}
