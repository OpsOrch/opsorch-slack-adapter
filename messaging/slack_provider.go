package messaging

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/opsorch/opsorch-core/messaging"
	"github.com/opsorch/opsorch-core/schema"
	"github.com/slack-go/slack"
)

const ProviderName = "slack"

// Config captures configuration for the Slack adapter.
type Config struct {
	Token string
}

type SlackProvider struct {
	client *slack.Client
}

func New(cfg map[string]any) (messaging.Provider, error) {
	parsed := parseConfig(cfg)
	if parsed.Token == "" {
		return nil, fmt.Errorf("slack token is required")
	}
	return &SlackProvider{
		client: slack.New(parsed.Token),
	}, nil
}

func init() {
	_ = messaging.RegisterProvider(ProviderName, New)
}

func (p *SlackProvider) Send(ctx context.Context, msg schema.Message) (schema.MessageResult, error) {
	var blocks []slack.Block

	// Map generic blocks to Slack blocks
	for _, b := range msg.Blocks {
		switch b.Type {
		case schema.BlockTypeHeader:
			blocks = append(blocks, slack.NewHeaderBlock(
				slack.NewTextBlockObject(slack.PlainTextType, b.Text, false, false),
			))
		case schema.BlockTypeSection:
			if len(b.Fields) > 0 {
				var fields []*slack.TextBlockObject
				for k, v := range b.Fields {
					// Format: *Key*\nValue
					text := fmt.Sprintf("*%s*\n%s", k, convertLinks(v))
					fields = append(fields, slack.NewTextBlockObject(slack.MarkdownType, text, false, false))
				}
				blocks = append(blocks, slack.NewSectionBlock(nil, fields, nil))
			} else {
				blocks = append(blocks, slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, convertLinks(b.Text), false, false),
					nil, nil,
				))
			}
		case schema.BlockTypeDivider:
			blocks = append(blocks, slack.NewDividerBlock())
		}
	}

	// Fallback to Body if no blocks provided
	if len(blocks) == 0 && msg.Body != "" {
		blocks = append(blocks, slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, convertLinks(msg.Body), false, false),
			nil, nil,
		))
	}

	opts := []slack.MsgOption{
		slack.MsgOptionBlocks(blocks...),
	}

	if msg.ThreadRef != "" {
		opts = append(opts, slack.MsgOptionTS(msg.ThreadRef))
	}

	channelID, timestamp, err := p.client.PostMessageContext(ctx, msg.Channel, opts...)
	if err != nil {
		return schema.MessageResult{}, err
	}

	return schema.MessageResult{
		ID:      timestamp,
		Channel: channelID,
		SentAt:  time.Now(),
	}, nil
}

func parseConfig(cfg map[string]any) Config {
	out := Config{}
	if v, ok := cfg["token"].(string); ok {
		out.Token = v
	}
	return out
}

// convertLinks converts [text](url) to <url|text>
var linkRegex = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

func convertLinks(input string) string {
	return linkRegex.ReplaceAllString(input, "<$2|$1>")
}
