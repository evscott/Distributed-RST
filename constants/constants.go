package constants

type Intent string

const (
	IntentPing         Intent = "ping"
	IntentSendPosition Intent = "send position"
	IntentSendGo	   Intent = "send go"
	IntentSendBack      Intent = "send back"
)