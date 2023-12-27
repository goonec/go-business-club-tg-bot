package openai

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestOpenAI_ResponseGPT(t *testing.T) {
	//cfg, err := config.New()
	//if err != nil {
	//	require.Nil(t, err)
	//}

	opneai := NewOpenAIConnect("")

	tests := []struct {
		text    string
		wantErr error
	}{
		{
			text:    "Напиши как работают горутины в go",
			wantErr: nil,
		},
		{
			text:    "Напиши все законы Ньюотона",
			wantErr: nil,
		},
		{
			text:    "Кто такое Уткин",
			wantErr: nil,
		},
		{
			text:    "Сколько ты выдержишь запросов?",
			wantErr: nil,
		},
		{
			text:    "Напиши бизнес схему по развитию цветочного бизнеса",
			wantErr: nil,
		},
		{
			text:    "Как дела?",
			wantErr: nil,
		},
	}

	var wg sync.WaitGroup

	for _, tt := range tests {
		wg.Add(1)
		go func(msg string) {
			defer wg.Done()
			t.Run(msg, func(t *testing.T) {
				response, err := opneai.ResponseGPT(msg)
				assert.Equal(t, err, tt.wantErr)
				if response == "" {
					t.Error("Message is empty")
				}
				t.Log(response)
			})
		}(tt.text)
	}

	wg.Wait()
}
