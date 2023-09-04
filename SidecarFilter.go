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
)

type SidecarFilter struct {
	Function TriFunction
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

// TODO solo copiar lo que hay dentro del propio nombre del body
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
	filter.Function(req, res, chain)
}

func (filter *SidecarFilter) Listen() {
	port := os.Getenv("PPORT")
	http.Handle("/", filter)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err)
	}
}
