package commonlog

const (
	MESSAGE = "_message"
	SCOPE   = "_scope"
	FILE    = "_file"
	LINE    = "_line"
)

//
// Message
//

// The entry point into the CommonLog API.
//
// Also see [Logger] as a more familiar, alternative API.
type Message interface {
	// Sets a value on the message and returns the same message
	// object.
	//
	// These keys are often specially supported:
	//
	// "_message": the base text of the message
	// "_scope": the scope of the message
	// "_file": filename in which the message was created
	// "_line": line number in the "_file"
	Set(key string, value any) Message

	// Sends the message to the backend.
	Send()
}
