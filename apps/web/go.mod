// This boundary keeps the root Go module from traversing Node dependencies
// under apps/web during `go test ./...`.
module github.com/ceheng-io/stock-go/apps/web

go 1.22
