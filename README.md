# grafana
- Command to start Docker container
 `docker run --name grafana --rm -v "$PWD/data:/etc/grafana" -p 3000:3000 grafana/grafana`
- Command to start Go server
  `go run main.go`
  - Server will start at port `8090` and hit `localhost:8090/grafana/login`
