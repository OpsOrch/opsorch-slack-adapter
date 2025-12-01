//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/opsorch/opsorch-core/schema"
	"github.com/opsorch/opsorch-slack-adapter/messaging"
)

func main() {
	token := os.Getenv("SLACK_TOKEN")
	channel := os.Getenv("SLACK_CHANNEL")

	if token == "" || channel == "" {
		fmt.Println("⚠️  Skipping integration tests: SLACK_TOKEN and SLACK_CHANNEL env vars are required.")
		return
	}

	fmt.Println("=================================")
	fmt.Println("Slack Adapter Integration Test")
	fmt.Println("=================================")

	ctx := context.Background()
	config := map[string]any{
		"token": token,
	}

	provider, err := messaging.New(config)
	if err != nil {
		log.Fatalf("❌ Failed to create provider: %v", err)
	}

	// Test 1: Send a simple message
	fmt.Println("\n=== Test 1: Send Simple Message ===")
	res, err := provider.Send(ctx, schema.Message{
		Channel: channel,
		Body:    "Integration test: Simple message from OpsOrch",
	})
	if err != nil {
		log.Printf("❌ Failed to send simple message: %v", err)
	} else {
		fmt.Printf("✅ Sent simple message. ID: %s, Timestamp: %s\n", res.ID, res.SentAt)
	}

	// Test 2: Send a rich message with blocks
	fmt.Println("\n=== Test 2: Send Rich Message (Blocks) ===")
	res, err = provider.Send(ctx, schema.Message{
		Channel: channel,
		Blocks: []schema.Block{
			{Type: schema.BlockTypeHeader, Text: "Integration Test Alert"},
			{Type: schema.BlockTypeSection, Text: "This is a *rich* message with [links](https://opsorch.com)."},
			{Type: schema.BlockTypeSection, Fields: map[string]string{
				"Environment": "Integration",
				"Status":      "Testing",
				"Link":        "[Click Me](https://google.com)",
			}},
			{Type: schema.BlockTypeDivider},
			{Type: schema.BlockTypeSection, Text: "End of test message."},
		},
	})
	if err != nil {
		log.Printf("❌ Failed to send rich message: %v", err)
	} else {
		fmt.Printf("✅ Sent rich message. ID: %s, Timestamp: %s\n", res.ID, res.SentAt)
	}
}
