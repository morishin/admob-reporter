package reporter

type AdmobReport []struct {
	Header AdmobReportHeader `json:"header,omitempty"`
	Row    AdmobReportRow    `json:"row,omitempty"`
	Footer AdmobReportFooter `json:"footer,omitempty"`
}

type AdmobReportHeader struct {
	DateRange struct {
		StartDate struct {
			Year  int `json:"year"`
			Month int `json:"month"`
			Day   int `json:"day"`
		} `json:"startDate"`
		EndDate struct {
			Year  int `json:"year"`
			Month int `json:"month"`
			Day   int `json:"day"`
		} `json:"endDate"`
	} `json:"dateRange"`
	LocalizationSettings struct {
		CurrencyCode string `json:"currencyCode"`
	} `json:"localizationSettings"`
}

type AdmobReportRow struct {
	DimensionValues struct {
		Date struct {
			Value string `json:"value"`
		} `json:"DATE"`
	} `json:"dimensionValues"`
	MetricValues struct {
		EstimatedEarnings struct {
			MicrosValue string `json:"microsValue"`
		} `json:"ESTIMATED_EARNINGS"`
	} `json:"metricValues"`
}

type AdmobReportFooter struct {
	MatchingRowCount string `json:"matchingRowCount"`
}

type Earnings struct {
	Yesterday int
	Today     int
	LastMonth int
	ThisMonth int
}

type OAuth2ClientSecret struct {
	Web struct {
		ClientID                string `json:"client_id"`
		ProjectID               string `json:"project_id"`
		AuthURI                 string `json:"auth_uri"`
		TokenURI                string `json:"token_uri"`
		AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
		ClientSecret            string `json:"client_secret"`
	} `json:"web"`
}
