package main

import "dist_sys/internal/server"

func main() {
	server := server.NewHTTPServer(":8080")
	server.ListenAndServe()
}
