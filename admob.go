package reporter

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jinzhu/now"
	"golang.org/x/oauth2"
	"google.golang.org/api/admob/v1"
	"google.golang.org/api/googleapi"
)

//go:embed oauth_client_secret.json
var clientSecretJson []byte

func loadOAuthConfig() *oauth2.Config {
	var clientSecret OAuth2ClientSecret
	if err := json.Unmarshal(clientSecretJson, &clientSecret); err != nil {
		panic(err)
	}
	return &oauth2.Config{
		ClientID:     clientSecret.Web.ClientID,
		ClientSecret: clientSecret.Web.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/admob.readonly",
			"https://www.googleapis.com/auth/admob.report",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  clientSecret.Web.AuthURI,
			TokenURL: clientSecret.Web.TokenURI,
		},
	}
}

func getReports(publisherId string, refreshToken string) (*AdmobReport, *AdmobReport) {
	ctx := context.Background()
	token := &oauth2.Token{
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(-1), // To always use the refresh token
	}
	config := loadOAuthConfig()
	tokenSource := config.TokenSource(ctx, token)
	client := oauth2.NewClient(ctx, tokenSource)

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	dailyReportReqBody := admob.GenerateNetworkReportRequest{
		ReportSpec: &admob.NetworkReportSpec{
			Dimensions: []string{"DATE"},
			Metrics:    []string{"ESTIMATED_EARNINGS"},
			DateRange: &admob.DateRange{
				StartDate: &admob.Date{
					Day:   int64(yesterday.Day()),
					Month: int64(yesterday.Month()),
					Year:  int64(yesterday.Year()),
				},
				EndDate: &admob.Date{
					Day:   int64(today.Day()),
					Month: int64(today.Month()),
					Year:  int64(today.Year()),
				},
			},
		},
	}

	firstDayOfLastMonth := now.BeginningOfMonth().AddDate(0, -1, 0)
	monthlyReportReqBody := admob.GenerateNetworkReportRequest{
		ReportSpec: &admob.NetworkReportSpec{
			Dimensions: []string{"MONTH"},
			Metrics:    []string{"ESTIMATED_EARNINGS"},
			DateRange: &admob.DateRange{
				StartDate: &admob.Date{
					Day:   int64(firstDayOfLastMonth.Day()),
					Month: int64(firstDayOfLastMonth.Month()),
					Year:  int64(firstDayOfLastMonth.Year()),
				},
				EndDate: &admob.Date{
					Day:   int64(today.Day()),
					Month: int64(today.Month()),
					Year:  int64(today.Year()),
				},
			},
		},
	}

	var dailyReport *AdmobReport
	var monthlyReport *AdmobReport
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		res := doRequest(client, publisherId, dailyReportReqBody)
		dailyReport = &res
	}()
	go func() {
		defer wg.Done()
		res := doRequest(client, publisherId, monthlyReportReqBody)
		monthlyReport = &res
	}()
	wg.Wait()

	return dailyReport, monthlyReport
}

func reportToEarnings(dailyReport *AdmobReport, monthlyReport *AdmobReport) Earnings {
	lastMonthEarnings, _ := strconv.Atoi((*monthlyReport)[0].Row.MetricValues.EstimatedEarnings.MicrosValue)
	var thisMonthEarnings = 0
	if len(*monthlyReport) >= 2 {
		thisMonthEarnings, _ = strconv.Atoi((*monthlyReport)[1].Row.MetricValues.EstimatedEarnings.MicrosValue)
	}
	yesterdayEarnings, _ := strconv.Atoi((*dailyReport)[0].Row.MetricValues.EstimatedEarnings.MicrosValue)
	var todayEarnings = 0
	if len(*dailyReport) >= 2 {
		todayEarnings, _ = strconv.Atoi((*dailyReport)[1].Row.MetricValues.EstimatedEarnings.MicrosValue)
	}
	earnings := Earnings{
		LastMonth: int(lastMonthEarnings / 1000000.0),
		ThisMonth: int(thisMonthEarnings / 1000000.0),
		Yesterday: int(yesterdayEarnings / 1000000.0),
		Today:     int(todayEarnings / 1000000.0),
	}
	return earnings
}

func doRequest(client *http.Client, publisherId string, reqBody admob.GenerateNetworkReportRequest) AdmobReport {
	body, _ := reqBody.MarshalJSON()
	url := fmt.Sprintf("https://admob.googleapis.com/v1/accounts/%s/networkReport:generate", publisherId)
	res, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if err := googleapi.CheckResponse(res); err != nil {
		panic(err)
	}

	var decoded AdmobReport
	json.NewDecoder(res.Body).Decode(&decoded)
	rows := decoded[1 : len(decoded)-1]
	return rows
}
