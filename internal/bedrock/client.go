// Package bedrock provides a client for AWS Bedrock AI inference.
package bedrock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// Analyser generates AI analysis text given a prompt.
type Analyser interface {
	Analyse(ctx context.Context, prompt string) (string, error)
}

// Client calls AWS Bedrock to generate analysis text.
type Client struct {
	runtimeClient *bedrockruntime.Client
	modelID       string
}

// NewClient constructs a Client for the given region and model.
// AWS credentials are resolved from the default credential chain (EC2 instance profile in prod).
func NewClient(region, modelID string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	return &Client{
		runtimeClient: bedrockruntime.NewFromConfig(cfg),
		modelID:       modelID,
	}, nil
}

type bedrockRequest struct {
	AnthropicVersion string    `json:"anthropic_version"`
	MaxTokens        int       `json:"max_tokens"`
	Messages         []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type bedrockResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

// Analyse sends a prompt to Bedrock and returns the generated text.
func (c *Client) Analyse(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(bedrockRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        512,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal bedrock request: %w", err)
	}

	out, err := c.runtimeClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.modelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return "", fmt.Errorf("invoke bedrock model: %w", err)
	}

	var resp bedrockResponse
	if err := json.Unmarshal(out.Body, &resp); err != nil {
		return "", fmt.Errorf("unmarshal bedrock response: %w", err)
	}
	if len(resp.Content) == 0 {
		return "", fmt.Errorf("bedrock returned empty content")
	}
	return resp.Content[0].Text, nil
}
