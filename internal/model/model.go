package model

type ServiceMMS map[int]MMSData
type ServiceSMS map[int]*SMSData

type ResultSetT struct {
	SMS       [][]SMSData            `json:"sms"`
	MMS       [][]MMSData            `json:"mms"`
	VoiceCall []VoiceCallData        `json:"voice_call"`
	Email     map[string][]EmailData `json:"email"`
	Billing   BiliingData            `json:"billing"`
	Support   []string               `json:"suppurt"`
	Incidents []IncidentData         `json:"incidents"`
}

type SMSData struct {
	Country      string
	Bandwidth    string
	ResponseTime string
	Provider     string
}

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

type VoiceCallData struct {
	Country             string
	Bandwidth           string
	ResponseTime        string
	Provider            string
	ConnectionStability float32
	TTFB                int
	VoicePurity         int
	MediaOfCallsTime    int
}

type EmailData struct {
	Country      string
	Provider     string
	DeliveryTime int
}

type BiliingData struct {
	CreateCustomer bool
	Purchase       bool
	Payout         bool
	Recurring      bool
	FraudControl   bool
	CheckoutPage   bool
}

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}
