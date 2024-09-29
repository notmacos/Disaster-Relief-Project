package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	if os.Args[1] == "message" {
		// Check if command-line arguments are provided
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run eventRecommendations.go message <message>")
			os.Exit(1)
		}
		fmt.Println(checkMessage(os.Args[2]))
	} else if os.Args[1] == "summary" {
		if len(os.Args) < 7 {
			fmt.Println("Usage: go run eventRecommendations.go summary <destination> <weather> <special_notes> <event_type> <user type>")
			os.Exit(1)
		}
		// Parse command-line arguments
		dataPoints := map[string]string{
			"destination":   os.Args[2],
			"weather":       os.Args[3],
			"special_notes": os.Args[4],
			"event type":    os.Args[5],
		}
		userType := os.Args[6]

		generateRecommendations(dataPoints, userType)
	}

}

func generateRecommendations(dataPoints map[string]string, userType string) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")

	// Configure the model for more focused and safe responses
	model.SetTemperature(0.2)
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	// Construct the prompt
	promptParts := []string{
		"Based on the following information, provide recommendations for what the user should bring, do, and be cautious about:",
	}
	for key, value := range dataPoints {
		promptParts = append(promptParts, fmt.Sprintf("%s: %s", key, value))
	}
	if userType == "volunteer" {
		promptParts = append(promptParts, "\nPlease format your response in three sections: 'Items to Bring', 'Recommended Actions', and 'Cautions'. Make sure to keep the response concise and to the point.")
	} else if userType == "official" {
		promptParts = append(promptParts, "\nBased on the following information, provide a concise summary of the incident, including key details such as the nature of the emergency, location, affected individuals, and any specific challenges or concerns. Do not tell the user what to do, just provide the information.")
	}

	prompt := strings.Join(promptParts, "\n")

	response, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
	}

	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			fmt.Println(part)
		}
	}
}

func checkMessage(message string) string {
	// Will check if message is SFW or NSFW
	// If NSFW, will return "[message restricted]"
	// If SFW, will return true
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")

	// Configure the model for content safety check
	model.SetTemperature(0.0)
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}

	prompt := fmt.Sprintf("Is the following message SFW (Safe For Work) or NSFW (Not Safe For Work)? Respond with only 'SFW' or 'NSFW'.\n\nMessage: %s", message)

	response, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Failed to generate content: %v", err)
		return "[message restricted]"
	}

	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		if part, ok := response.Candidates[0].Content.Parts[0].(genai.Text); ok {
			result := string(part)
			if strings.TrimSpace(strings.ToUpper(result)) == "NSFW" {
				return "[message restricted]"
			}
		}
	}

	return message
}
