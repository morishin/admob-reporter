package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestCdkStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)
	setUpEnv(t)

	// WHEN
	stack := NewCdkStack(app, "MyStack", nil)

	// THEN
	bytes, err := json.Marshal(app.Synth(nil).GetStackArtifact(stack.ArtifactId()).Template())
	if err != nil {
		t.Error(err)
	}

	template := gjson.ParseBytes(bytes)
	functionName := template.Get("Resources.admobreporterfunctionE6B4D67C.Properties.FunctionName").String()
	cronExpression := template.Get("Resources.admobreporterrule538BDEC9.Properties.ScheduleExpression").String()
	env := template.Get("Resources.admobreporterfunctionE6B4D67C.Properties.Environment.Variables").Map()
	assert.True(t, strings.HasPrefix(env["ADMOB_PUBLISHER_ID"].String(), "pub-"))
	assert.Equal(t, "admob-reporter-function", functionName)
	assert.Equal(t, "cron(0 3,15 * * ? *)", cronExpression)
}

func setUpEnv(t *testing.T) {
	t.Setenv("ADMOB_PUBLISHER_ID", "pub-xxxxxxxxxxxxxxxxx")
	t.Setenv("ADMOB_OAUTH2_REFRESH_TOKEN", "XXXXXXXXXX")
	t.Setenv("SLACK_WEBHOOK_URL", "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX")
	t.Setenv("TZ", "Asia/Tokyo")
}
