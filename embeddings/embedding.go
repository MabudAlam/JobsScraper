package embeddings

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"jobscraper/common"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/operations"
)

const (
	EmbeddingModel    = "openai/text-embedding-3-small"
	BatchSize         = 2048
	MaxRetries        = 3
	MaxCharsPerText   = 6000
	EmbeddingDims     = 1536
	OpenRouterBaseURL = "https://openrouter.ai/api/v1"
)

type EmbeddingService struct {
	client *openrouter.OpenRouter
}

func NewEmbeddingService(apiKey string) *EmbeddingService {
	client := openrouter.New(
		openrouter.WithSecurity(apiKey),
	)
	return &EmbeddingService{client: client}
}

func (s *EmbeddingService) GenerateJobText(job *common.JobPayload) string {
	var parts []string

	if job.JobName != "" {
		parts = append(parts, job.JobName)
	}
	if job.Description != "" {
		parts = append(parts, job.Description)
	}

	return strings.Join(parts, " ")
}

type EmbedResult struct {
	Index  int
	Vector []float32
}

func (s *EmbeddingService) EmbedTexts(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	allEmbeddings := make([][]float32, len(texts))

	for batchStart := 0; batchStart < len(texts); batchStart += BatchSize {
		batchEnd := batchStart + BatchSize
		if batchEnd > len(texts) {
			batchEnd = len(texts)
		}
		batch := texts[batchStart:batchEnd]

		nonEmptyIndices := make([]int, 0)
		textsToEmbed := make([]string, 0)

		for i, text := range batch {
			truncated := text
			if len(truncated) > MaxCharsPerText {
				truncated = truncated[:MaxCharsPerText]
			}

			if truncated != "" && strings.TrimSpace(truncated) != "" {
				nonEmptyIndices = append(nonEmptyIndices, i)
				textsToEmbed = append(textsToEmbed, truncated)
			} else {
				allEmbeddings[batchStart+i] = make([]float32, EmbeddingDims)
			}
		}

		totalChars := 0
		for _, t := range textsToEmbed {
			totalChars += len(t)
		}

		log.Printf("[Embedding] Batch %d-%d of %d texts (%d non-empty, ~%d chars total)",
			batchStart, batchEnd, len(texts), len(textsToEmbed), totalChars)

		if len(textsToEmbed) == 0 {
			log.Printf("[Embedding] WARNING: Batch %d-%d has no non-empty texts; using zero vectors",
				batchStart, batchEnd)
			continue
		}

		vectors, err := s.embedWithRetry(textsToEmbed)
		if err != nil {
			return nil, fmt.Errorf("failed to embed batch %d-%d: %w", batchStart, batchEnd, err)
		}

		for outIdx, inIdx := range nonEmptyIndices {
			allEmbeddings[batchStart+inIdx] = vectors[outIdx]
		}
	}

	return allEmbeddings, nil
}

func (s *EmbeddingService) embedWithRetry(texts []string) ([][]float32, error) {
	var lastErr error

	for attempt := 0; attempt < MaxRetries; attempt++ {
		if attempt > 0 {
			wait := time.Duration(1<<attempt) * time.Second
			log.Printf("[Embedding] Retry attempt %d/%d, waiting %v", attempt+1, MaxRetries, wait)
			time.Sleep(wait)
		}

		vectors, err := s.embedBatch(texts)
		if err == nil {
			return vectors, nil
		}
		lastErr = err
		log.Printf("[Embedding] Attempt %d/%d failed: %v", attempt+1, MaxRetries, err)
	}

	return nil, fmt.Errorf("embedding failed after %d attempts: %w", MaxRetries, lastErr)
}

func (s *EmbeddingService) embedBatch(texts []string) ([][]float32, error) {
	input := operations.CreateInputUnionArrayOfStr(texts)

	req := operations.CreateEmbeddingsRequest{
		Input:          input,
		Model:          EmbeddingModel,
		EncodingFormat: operations.EncodingFormatFloat.ToPointer(),
	}

	resp, err := s.client.Embeddings.Generate(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("embedding request failed: %w", err)
	}

	if resp.CreateEmbeddingsResponseBody != nil {
		data := resp.CreateEmbeddingsResponseBody
		if len(data.Data) == 0 {
			return nil, fmt.Errorf("OpenRouter returned empty data for batch of %d texts", len(texts))
		}

		embeddings := make([][]float32, 0, len(data.Data))
		for _, d := range data.Data {
			if emb := d.Embedding; emb.ArrayOfNumber != nil {
				floatEmbed := make([]float32, len(emb.ArrayOfNumber))
				for i, v := range emb.ArrayOfNumber {
					floatEmbed[i] = float32(v)
				}
				embeddings = append(embeddings, floatEmbed)
			}
		}
		return embeddings, nil
	}

	if resp.Str != nil {
		return nil, fmt.Errorf("unexpected response type: %s", *resp.Str)
	}

	return nil, fmt.Errorf("unknown response type")
}

func (s *EmbeddingService) EmbedJobs(ctx context.Context, jobs []*common.JobPayload) ([]*common.JobPayload, error) {
	if len(jobs) == 0 {
		return jobs, nil
	}

	texts := make([]string, len(jobs))
	for i, job := range jobs {
		texts[i] = s.GenerateJobText(job)
	}

	embeddings, err := s.EmbedTexts(texts)
	if err != nil {
		return nil, err
	}

	for i, job := range jobs {
		if i < len(embeddings) {
			job.Embedding = embeddings[i]
		}
	}

	return jobs, nil
}

func EmbedJobTexts(texts []string) ([][]float32, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY not set")
	}

	svc := NewEmbeddingService(apiKey)
	return svc.EmbedTexts(texts)
}
