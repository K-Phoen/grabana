package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/alertmanager"
	"github.com/K-Phoen/grabana/alertmanager/email"
	"github.com/K-Phoen/grabana/alertmanager/opsgenie"
	"github.com/K-Phoen/grabana/alertmanager/slack"
	"github.com/K-Phoen/grabana/alertmanager/webhook"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Usage: go run main.go http://grafana-host:3000 api-key-string-here\n")
		os.Exit(1)
	}

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, os.Args[1], grabana.WithAPIToken(os.Args[2]))

	manager := alertmanager.New(
		alertmanager.ContactPoints(
			alertmanager.ContactPoint(
				"Platform",
				email.To([]string{"joe@lafrite"}),
				opsgenie.With("https://api.eu.opsgenie.com/v2/alerts", "some API key", opsgenie.AutoClose(), opsgenie.OverridePriority()),
				slack.Webhook("https://api.slack.com/messaging/webhooks"),
			),
			alertmanager.ContactPoint(
				"Core Exp",
				email.To([]string{"core@exp"}, email.Single()),
				webhook.Call("http://example.com"),
			),
		),
		alertmanager.Routing(
			alertmanager.Policy("Platform", alertmanager.TagEq("owner", "platform")),
		),
		alertmanager.DefaultContactPoint("Core Exp"),
	)

	if err := client.ConfigureAlertManager(ctx, manager); err != nil {
		fmt.Printf("Could not configure alerting: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("The deed is done.")
}
