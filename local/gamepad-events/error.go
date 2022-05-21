package gamepad_events

type errorCode uint32
const (
	// Invalid JSON in load operation
	ErrBadJson errorCode = iota
	// Invalid name in the JSON
	ErrJsonInvalidName
	// Invalid GUID in the JSON
	ErrJsonInvalidGuid
	// Invalid input in the JSON
	ErrJsonInvalidInput
	// Couldn't decode the input into its type
	ErrJsonBadInput
	// Missing name in the JSON
	ErrJsonMissingName
	// Missing GUID in the JSON
	ErrJsonMissingGuid
	// There's no input in the JSON
	ErrJsonInputEmpty
)

func (e errorCode) Error() string {
	switch e {
	case ErrBadJson:
		return "Invalid JSON in load operation"
	case ErrJsonInvalidName:
		return "Invalid name in the JSON"
	case ErrJsonInvalidGuid:
		return "Invalid GUID in the JSON"
	case ErrJsonInvalidInput:
		return "Invalid input in the JSON"
	case ErrJsonBadInput:
		return "Couldn't decode the input into its type"
	case ErrJsonMissingName:
		return "Missing name in the JSON"
	case ErrJsonMissingGuid:
		return "Missing GUID in the JSON"
	case ErrJsonInputEmpty:
		return "There's no input in the JSON"
	default:
		return "Unknown error"
	}
}
