package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/decoder"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprint(os.Stderr, "Usage: go run -mod=vendor main.go file http://grafana-host:3000 api-key\n")
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %s\n", err)
		os.Exit(1)
	}

	dashboard, err := decoder.UnmarshalYAML(bytes.NewBuffer(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse file: %s\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, os.Args[2], grabana.WithAPIToken(os.Args[3]))

	// create the folder holding the dashboard for the service
	folder, err := client.GetFolderByTitle(ctx, "Test Folder")
	if err != nil && err != grabana.ErrFolderNotFound {
		fmt.Printf("Could not create folder: %s\n", err)
		os.Exit(1)
	}
	if folder == nil {
		folder, err = client.CreateFolder(ctx, "Test Folder")
		if err != nil {
			fmt.Printf("Could not create folder: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Folder created (id: %d, uid: %s)\n", folder.ID, folder.UID)
	}

	if _, err := client.UpsertDashboard(ctx, folder, dashboard); err != nil {
		fmt.Printf("Could not create dashboard: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("The deed is done.")
}
