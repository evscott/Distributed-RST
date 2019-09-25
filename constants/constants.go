package constants

const (
	StdIn  string = "<< %v: "
	StdOut string = ">> %v %v\n"
)

type Intent string

const (
	IntentPing         Intent = "ping"
	IntentSendPosition Intent = "send position"

)