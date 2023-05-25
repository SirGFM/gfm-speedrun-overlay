package brat_handler

// Errors returned by this package.
type errorCode uint

const (
	// Failed to decode the JSON input
	BadJSONInput errorCode = iota
	// Didn't find any resource associated with the given path
	ResourceNotFound
)

// Implement the `error` interface for `errorCode`.
func (e errorCode) Error() string {
	switch e {
	case BadJSONInput:
		return "brat-handler: Failed to decode the JSON input"
	case ResourceNotFound:
		return "brat-handler: Didn't find any resource associated with the given path"
	default:
		return "brat-handler: Unknown"
	}
}
