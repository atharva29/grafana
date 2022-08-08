package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func HandleHelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello World")
		sendResponse(r.Context(), w, http.StatusOK, map[string]interface{}{"message": "hello world"})
	}
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// read username from request body &
		// authenticate the user by verifying in database or any other way
		// pass the role of the user
		// "admin" role will have all privileges
		// rest roles will have access as per authorization policy
		user := "admin"
		req.Header.Set("X-WEBAUTH-USER", user)
	}
	return proxy, nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

// sendResponse sends an http response
func sendResponse(ctx context.Context, w http.ResponseWriter, statusCode int, body interface{}) error {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(body)
}

// CreateCorsObject creates a cors object with the required config
func createCorsObject() *cors.Cors {
	return cors.New(cors.Options{
		AllowCredentials: true,
		AllowOriginFunc: func(s string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		ExposedHeaders: []string{"Authorization", "Content-Type"},
	})
}

func main() {
	router := mux.NewRouter()
	proxy, err := NewProxy("http://localhost:3000")
	if err != nil {
		panic(err)
	}

	publicRoutes := router.Methods(http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete).Subrouter()
	publicRoutes.PathPrefix("/grafana/").HandlerFunc(ProxyRequestHandler(proxy))
	publicRoutes.Methods(http.MethodGet).Path("/hello").HandlerFunc(HandleHelloWorld())
	fmt.Println("HEllo world")
	corsObj := createCorsObject()
	Handler := corsObj.Handler(router)

	fmt.Println("Starting http server on port: " + strconv.Itoa(8090))
	http.ListenAndServe(":"+strconv.Itoa(8090), Handler)
}
