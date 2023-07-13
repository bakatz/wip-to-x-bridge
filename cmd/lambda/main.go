package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	lib_wip "github.com/bakatz/wip-to-twitter-bridge/lib/wip"
	"github.com/dghubble/oauth1"
	twitter "github.com/g8rswimmer/go-twitter/v2"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Response struct {
	Message         string `json:"message"`
	Code            string `json:"code,omitempty"`
	NumTodosTweeted int    `json:"num_todos_tweeted"`
}

const (
	LOOKBACK_WINDOW_MINUTES       = 60
	SUCCESS_MESSAGE               = "Function finished without errors"
	CONNECTION_TIMEOUT_DURATION   = 5 * time.Second
	CONTENT_TYPE_APPLICATION_JSON = "application/json"
	WIP_API_GRAPHQL_QUERY         = `{
		viewer {
			projects {
				id
				name
				website_url
				todos(completed:true, orderBy: { completedAt:desc }) {
					id
					body
					completed_at
					attachments {
						url
					}
				}
			}
		}
	}
	`
)

// Dummy auth struct just to satisfy the API
type authorize struct{}

func (a authorize) Add(req *http.Request) {}

func makeAndLogErrorResponse(message string, code string, logger *zap.Logger) Response {
	response := Response{Message: message, Code: code}
	logger.Sugar().Error("Response: ", response)
	return response
}

func Handler(ctx context.Context) (Response, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Get all the secrets we need
	wipAPIKey := os.Getenv("WIP_API_KEY")
	twitterAPIKey := os.Getenv("TWITTER_API_KEY")
	twitterAPIKeySecret := os.Getenv("TWITTER_API_KEY_SECRET")
	twitterAccessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	twitterAccessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	if wipAPIKey == "" || twitterAPIKey == "" || twitterAPIKeySecret == "" || twitterAccessToken == "" || twitterAccessTokenSecret == "" {
		return makeAndLogErrorResponse("Cannot start the function because some of the required evars are missing, set them and run the function again", "missing_evars", logger), nil
	}

	// Get all of the completed todos from wip.co
	graphQLRequestBytes, err := json.Marshal(map[string]string{
		"query": WIP_API_GRAPHQL_QUERY,
	})
	if err != nil {
		return makeAndLogErrorResponse("Error encoding WIP request body", "wip_request_encode_error", logger), err
	}
	req, err := http.NewRequest(http.MethodPost, "https://wip.co/graphql", bytes.NewBuffer(graphQLRequestBytes))
	if err != nil {
		return makeAndLogErrorResponse("Error creating WIP request", "wip_request_error", logger), err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+wipAPIKey)
	wipHttpClient := &http.Client{Timeout: CONNECTION_TIMEOUT_DURATION}
	resp, err := wipHttpClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return makeAndLogErrorResponse("Error calling WIP's API; status code = "+strconv.Itoa(resp.StatusCode), "wip_api_error", logger), err
	}
	wipResponse := &lib_wip.WIPAPIResponse{}
	err = json.NewDecoder(resp.Body).Decode(wipResponse)

	if err != nil {
		return makeAndLogErrorResponse("Error decoding WIP's response body", "wip_response_decode_error", logger), err
	}

	// Send out a tweet for each of the completed todos
	oauth1Config := oauth1.NewConfig(twitterAPIKey, twitterAPIKeySecret)
	twitterHttpClient := oauth1Config.Client(oauth1.NoContext, &oauth1.Token{
		Token:       twitterAccessToken,
		TokenSecret: twitterAccessTokenSecret,
	})
	twitterHttpClient.Timeout = CONNECTION_TIMEOUT_DURATION
	twitterClient := &twitter.Client{
		Authorizer: authorize{},
		Client:     twitterHttpClient,
		Host:       "https://api.twitter.com",
	}

	startOfLookbackWindow := time.Now().UTC().Add(-LOOKBACK_WINDOW_MINUTES * time.Minute)
	numTodosTweeted := 0
	for _, project := range wipResponse.Data.Viewer.Projects {
		for _, todo := range project.Todos {
			// If this todo was completed more than an hour ago, don't bother tweeting about it because we've already covered it in a previous run (we run every hour to catch todos from the previous hour)
			if todo.CompletedAt.Before(startOfLookbackWindow) {
				continue
			}
			tweetMessage := "Just completed a task for " + project.Name + ": " + todo.Body
			if len(todo.Attachments) > 0 {
				tweetMessage += " " + todo.Attachments[0].URL
			}
			logger.Info("About to tweet this message: " + tweetMessage)
			_, err := twitterClient.CreateTweet(context.Background(), twitter.CreateTweetRequest{
				Text: tweetMessage,
			})
			if err != nil {
				return makeAndLogErrorResponse("Error creating a tweet", "twitter_create_tweet_error", logger), err
			}
			logger.Info("Tweet sent successfully")
			numTodosTweeted++
		}
	}

	// Return a success message
	logger.Info(SUCCESS_MESSAGE, zap.Int("num_todos_tweeted", numTodosTweeted))
	return Response{Message: SUCCESS_MESSAGE, NumTodosTweeted: numTodosTweeted}, nil
}

func main() {
	godotenv.Load("../../.env")
	if os.Getenv("RUN_WITHOUT_LAMBDA") == "true" {
		Handler(context.TODO())
	} else {
		lambda.Start(Handler)
	}
}
