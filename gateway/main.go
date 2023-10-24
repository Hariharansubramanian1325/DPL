package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

// func main() {
// 	router := mux.NewRouter()

// 	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
// 		rw.Header().Set("Content-Type", "application/json")

// 		json.NewEncoder(rw).Encode(map[string]string{"data": "Hello from Mux & mongoDB"})
// 	}).Methods("GET")
// 	// configs.ConnectDB()
// 	routes.UserRoute(router)
// 	// log.Fatal(http.ListenAndServe(":6000", router))

// 	const port = ":8080"
// 	fmt.Println("Server started successfully on", port)
// 	server := &http.Server{
// 		Handler: router,
// 		Addr:    port, // Port and host
// 	}

// 	err2 := server.ListenAndServe()
// 	if err2 != nil {
// 		log.Fatal("cannot start server:", err2)
// 		return
// 	}

// }
func reverseProxyHandler(targetURL string) http.Handler {
	target, _ := url.Parse(targetURL)
	return httputil.NewSingleHostReverseProxy(target)
}
func main() {
	router := mux.NewRouter()
	backendService1URL := "http://localhost:8083"
	backendService2URL := "http://localhost:8082"
	backendService3URL := "http://localhost:8084"
	// Create routes for the API gateway
	router.PathPrefix("/service1").Handler(reverseProxyHandler(backendService1URL))
	router.PathPrefix("/service2").Handler(reverseProxyHandler(backendService2URL))
	router.PathPrefix("/service3").Handler(reverseProxyHandler(backendService3URL))
	http.Handle("/", router)

	go func() {
		fmt.Println("Server started successfully on port 8080")
		http.ListenAndServe(":8080", router)
	}()
	select {}
}
