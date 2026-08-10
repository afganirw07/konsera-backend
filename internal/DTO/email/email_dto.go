package email

type LoginOTPEmailData struct {
	Device           string
	Location         string
	IPAddress        string
	Time             string
	VerifyURL        string
	OTPCode          string
	OTPExpiry        int
	SecureAccountURL string
	UnsubscribeURL   string
	PrivacyURL       string
	HelpCenterURL    string
}
