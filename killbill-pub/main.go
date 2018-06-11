package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	me := os.Args[0]

	if len(os.Args) < 4 {
		log.Fatalf("usage: %s projectID topicID billingAccountID", me)
	}

	projectID := os.Args[1]
	topicID := os.Args[2]
	billingAccountID := os.Args[3]

	publish(projectID, topicID, billingAccountID)
}

func publish(projectID, topicID, billingAccountID string) {

	log.Printf("publish: project=%s topic=%s billingAccount=%s", projectID, topicID, billingAccountID)

	ctx := context.Background()
	client, errNew := pubsub.NewClient(ctx, projectID)
	if errNew != nil {
		log.Printf("publish: failure creating pubsub client: %v", errNew)
		return
	}
	defer client.Close()

	topic := client.Topic(topicID)
	defer topic.Stop()

	data := `
{ 
"budgetDisplayName":"My Personal Budget",
"alertThresholdExceeded":0.9,
"costAmount":140.321,
"costIntervalStart":"2018-02-01T08:00:00Z",
"budgetAmount":152.557,
"budgetAmountType":"SPECIFIED_AMOUNT",
"currencyCode":"USD"
}
`

	m := pubsub.Message{
		Data:       []byte(data),
		Attributes: map[string]string{},
	}
	m.Attributes["billingAccountId"] = billingAccountID
	m.Attributes["schemaVersion"] = "1.0"

	r := topic.Publish(ctx, &m)

	id, errPub := r.Get(ctx)
	if errPub != nil {
		log.Printf("publish: failure publishing message: %v", errPub)
		return
	}

	log.Printf("publish: published message with a message ID: %s\n", id)
}
