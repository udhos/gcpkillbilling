package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

const version = "0.0"

var dry = true
var once = true

func env(me, name string, value *bool) {
	if v := os.Getenv(name); v != "" {
		val, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatalf("%s: refusing to run with bad env var %s=%s: %v", me, name, v, err)
		}
		*value = val
	}
}

func main() {
	me := os.Args[0]

	if len(os.Args) < 3 {
		log.Fatalf("usage: %s projectID subscriptionID", me)
	}

	projectID := os.Args[1]
	subscriptionID := os.Args[2]

	env(me, "DRY", &dry)
	env(me, "ONCE", &once)

	log.Printf("%s: version=%s runtime=%s DRY=%v ONCE=%v", me, version, runtime.Version(), dry, once)

	pull(projectID, subscriptionID)
}

func pull(projectID, subscriptionID string) {
	log.Printf("pull: start project=%s subscription=%s", projectID, subscriptionID)

	ctx := context.Background()
	client, errNew := pubsub.NewClient(ctx, projectID)
	if errNew != nil {
		log.Printf("pull: failure creating pubsub client: %v", errNew)
		return
	}
	defer client.Close()

	sub := client.Subscription(subscriptionID)

	// This program is expected to process and acknowledge messages in 30 seconds. If
	// not, the Pub/Sub API will assume the message is not acknowledged.
	sub.ReceiveSettings.MaxExtension = 30 * time.Second
	sub.ReceiveSettings.MaxOutstandingMessages = 5
	sub.ReceiveSettings.MaxOutstandingBytes = 10e6

	var lock sync.RWMutex

	cctx, cancel := context.WithCancel(ctx)
	if errPull := sub.Receive(cctx, func(c context.Context, m *pubsub.Message) {
		// NOTE: May be called concurrently; synchronize access to shared memory.
		lock.RLock()
		handleMessage(c, m, cancel)
		lock.RUnlock()
	}); errPull != nil {
		log.Printf("Receive() error: %v", errPull)
	}

	log.Printf("pull: exit project=%s subscription=%s", projectID, subscriptionID)
}

func handleMessage(ctx context.Context, m *pubsub.Message, cancel context.CancelFunc) {
	log.Printf("pull: ID=%s PublishTime=%s\n", m.ID, m.PublishTime)
	log.Printf("pull: ID=%s data = %q\n", m.ID, m.Data)
	log.Printf("pull: ID=%s attributes = %v\n", m.ID, m.Attributes)

	billingAccountID, found := m.Attributes["billingAccountId"]
	if !found {
		log.Printf("pull: ID=%s missing billingAccountId=[%s] in message attributes", m.ID, billingAccountID)
		return
	}

	if errKill := killbill(billingAccountID); errKill != nil {
		log.Printf("pull: ID=%s failure killing billing for account=%s: %v", m.ID, billingAccountID, errKill)
		return
	}

	log.Printf("pull: ID=%s DRY=%v removing from queue", m.ID, dry)
	if !dry {
		m.Ack() // message handled
	}

	log.Printf("pull: ID=%s ONCE=%v finishing", m.ID, once)
	if once {
		cancel() // request termination
	}
}
