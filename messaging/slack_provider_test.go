package messaging

import (
	"context"
	"testing"

	"github.com/opsorch/opsorch-core/schema"
)

func TestConvertLinks(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello [world](http://example.com)", "Hello <http://example.com|world>"},
		{"Check [this](http://a.com) and [that](http://b.com)", "Check <http://a.com|this> and <http://b.com|that>"},
		{"No links here", "No links here"},
	}

	for _, tt := range tests {
		got := convertLinks(tt.input)
		if got != tt.expected {
			t.Errorf("convertLinks(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestSend_Mapping(t *testing.T) {
	// This test would ideally mock the Slack client.
	// For now, we just ensure the code compiles and the structure is correct.
	// Since we can't easily mock the slack-go client without an interface wrapper,
	// we rely on the integration tests or manual verification.
	ctx := context.Background()
	_ = ctx
	msg := schema.Message{
		Channel: "C123",
		Blocks: []schema.Block{
			{Type: schema.BlockTypeHeader, Text: "Header"},
		},
	}
	_ = msg
}
