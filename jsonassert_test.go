package jsonassert_test

import (
	"fmt"
	"testing"

	"github.com/kinbiko/jsonassert"
)

type fakeT struct {
	receivedMessages []string
}

func (ft *fakeT) Errorf(format string, args ...interface{}) {
	ft.receivedMessages = append(ft.receivedMessages, fmt.Sprintf(format, args...))
}

// Should be able to make assertions against the String representation of a
// JSON payload
func TestAssertString(t *testing.T) {
	ft := new(fakeT)
	ja := jsonassert.New(ft)
	tt := []struct {
		payload       string
		assertedJSON  string
		args          []interface{}
		expAssertions []string
	}{
		{
			payload:      `{"check": "nope", "ok": "nah"}`,
			assertedJSON: `{"check": "%s", "ok": "yup"}`,
			args:         []interface{}{"works"},
			expAssertions: []string{
				`Expected key: "check" to have value "works" but was "nope"`,
				`Expected key: "ok" to have value "yup" but was "nah"`,
			},
		},
	}
	for _, tc := range tt {
		ja.AssertString(tc.payload, tc.assertedJSON, tc.args...)

		msgs := ft.receivedMessages
		if exp, got := len(tc.expAssertions), len(msgs); exp != got {
			t.Errorf("Expected %d error messages to be written, but there were %d", exp, got)
			t.Errorf("Expected the following messages:")
			for _, msg := range tc.expAssertions {
				t.Errorf(msg)
			}
			t.Errorf("Got the following messages:")
			for _, msg := range msgs {
				t.Errorf(msg)
			}
			return //Don't attempt the following assertions

		}

		// The order of the JSON does not matter, so have to do a double subset check
		// Combines the issues in the end in order to make deciphering the test failure easier to parse
		unexpectedAssertions := ""
		for _, got := range msgs {
			found := false
			for _, exp := range tc.expAssertions {
				if got == exp {
					found = true
				}
			}
			if !found {
				if unexpectedAssertions == "" {
					unexpectedAssertions = "Got unexpected assertion failure:"
				}
				unexpectedAssertions += "\n - " + got
			}
		}

		missingAssertions := ""
		for _, got := range tc.expAssertions {
			found := false
			for _, exp := range msgs {
				if got == exp {
					found = true
				}
			}
			if !found {
				if missingAssertions == "" {
					missingAssertions = "\nExpected assertion failure but was not found:"
				}
				missingAssertions += "\n - " + got
			}
		}

		if totalError := unexpectedAssertions + missingAssertions; totalError != "" {
			t.Errorf("Inconsistent assertions:\n%s", totalError)
		}
	}
}
