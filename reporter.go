package reporter

import (
	"os"

	"github.com/joho/godotenv"
)

func Run() {
	publisherId, refreshToken, slackWebhookUrl := loadEnv()
	dailyReport, monthlyReport := getReports(publisherId, refreshToken)
	earnings := reportToEarnings(dailyReport, monthlyReport)
	postSlack(slackWebhookUrl, &earnings)
}

func loadEnv() (string, string, string) {
	godotenv.Load()
	publisherId := os.Getenv("ADMOB_PUBLISHER_ID")
	refreshToken := os.Getenv("ADMOB_OAUTH2_REFRESH_TOKEN")
	slackWebhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
	return publisherId, refreshToken, slackWebhookUrl
}
