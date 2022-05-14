package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/datasource/prometheus"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Usage: go run main.go http://grafana-host:3000 api-key-string-here\n")
		os.Exit(1)
	}

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, os.Args[1], grabana.WithAPIToken(os.Args[2]))

	options := []prometheus.Option{
		prometheus.Default(),
		prometheus.SkipTLSVerify(),
	}

	datasource, err := prometheus.New("grabana-prometheus", "http://172.17.0.1:9090", options...)
	if err != nil {
		fmt.Printf("Could not build datasource: %s\n", err)
		os.Exit(1)
	}

	if err := client.UpsertDatasource(ctx, datasource); err != nil {
		fmt.Printf("Could not create datasource: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("The deed is done.")
}
