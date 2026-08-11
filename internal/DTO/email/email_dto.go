package email

type RegisterOTPEmailData struct {
	Device           string
	Location         string
	IPAddress        string
	Time             string
	VerifyURL        string
	OTPCode          int
	OTPExpiry        int
	SecureAccountURL string
	UnsubscribeURL   string
	PrivacyURL       string
	HelpCenterURL    string
}

type RegisterOTPEmailMeta struct {
	Device    string
	Location  string
	IPAddress string
	Time      string
}