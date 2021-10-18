package reporter

import (
	"fmt"

	"github.com/slack-go/slack"
)

func postSlack(webhookUrl string, earnings *Earnings) {
	formatPrice := func(price int) string {
		if price < 10000 {
			return fmt.Sprintf("￥%d", price)
		} else {
			return fmt.Sprintf("￥%.2f万", float32(price)/10000)
		}
	}

	fields := []slack.AttachmentField{
		{
			Title: "本日 (現時点まで)",
			Value: formatPrice(earnings.Today),
			Short: true,
		}, {
			Title: "昨日",
			Value: formatPrice(earnings.Yesterday),
			Short: true,
		}, {
			Title: "今月 (現時点まで)",
			Value: formatPrice(earnings.ThisMonth),
			Short: true,
		}, {
			Title: "先月",
			Value: formatPrice(earnings.LastMonth),
			Short: true,
		},
	}
	attachments := []slack.Attachment{{
		Color:      "#f2a600",
		AuthorName: "AdMob ネットワークでの見積もり収益額",
		Fields:     fields,
	}}
	message := slack.WebhookMessage{
		Username:    "admob-reporter",
		Attachments: attachments,
	}
	slack.PostWebhook(webhookUrl, &message)
}
