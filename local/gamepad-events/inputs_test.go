package gamepad_events

import (
	"testing"
	"encoding/json"
)

// Check that an InputName may properly be encoded into a JSON, and then
// decoded from it.
func TestInputNameToJson(t *testing.T) {
	for i := InputName(0); i < NumInputNames; i++ {
		dict := make(map[InputName]int)
		dict[i] = int(i)

		data, err := json.Marshal(dict)
		if err != nil {
			t.Fatalf("Couldn't encode '%s' to JSON: %+v", i.key(), err)
		}

		strMap := make(map[string]int)
		err = json.Unmarshal(data, &strMap)
		if err != nil {
			t.Fatalf("Couldn't decode '%s' from a JSON: %+v", i.key(), err)
		}

		if val, ok := strMap[i.key()]; len(strMap) != 1 || !ok || val != int(i) {
			t.Fatalf("Failed to decode the JSON into a string. Got '%+v', Wanted '%+v'",
				dict,
				strMap)
		}

		decodedDict := make(map[InputName]int)
		err = json.Unmarshal(data, &decodedDict)
		if err != nil {
			t.Fatalf("Couldn't decode '%s' from a JSON: %+v", i.key(), err)
		}

		if val, ok := decodedDict[i]; len(strMap) != 1 || !ok || val != int(i) {
			t.Fatalf("Failed to decode the JSON into an InputName. Got '%+v', Wanted '%+v'",
				dict,
				strMap)
		}
	}
}
