package ai

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func SpeechToText(filepath string) (*string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file")
	}

	ctx := context.Background()

	key := os.Getenv("OPEN_AI_KEY")
	client := openai.NewClient(option.WithAPIKey(key))
	transcription, err := client.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
		File:  openai.F[io.Reader](file),
		Model: openai.F(openai.AudioModelWhisper1),
	})
	if err != nil {
		return nil, err
	}

	return &transcription.Text, nil
}
