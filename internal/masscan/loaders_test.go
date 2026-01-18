package masscan

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeValue(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		data        string
		expectValue any
		expectError string
	}{
		{
			"empty",
			``,
			"",
			"",
		},
		{
			"slice (json)",
			`["one", "two", ""]`,
			[]string{"one", "two", ""},
			"",
		},
		{
			"slice (comma)",
			`one,two,,three,`,
			[]string{"one", "two", "three"},
			"",
		},
		{
			"slice (newline)",
			"one\ntwo\nthree\n",
			[]string{"one", "two", "three"},
			"",
		},
		{
			"slice (single)",
			`one`,
			[]string{"one"},
			"",
		},
		{
			"string (raw)",
			`one`,
			"one",
			"",
		},
		{
			"string (json)",
			`"one"`,
			"one",
			"",
		},
		{
			"invalid json",
			`[test`,
			[]string{},
			"invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				result any
				err    error
			)

			switch tc.expectValue.(type) {
			case string:
				result, err = decodeValue[string](t.Context(), "test", "test-ref", []byte(tc.data))
			case []string:
				result, err = decodeValue[[]string](t.Context(), "test", "test-ref", []byte(tc.data))
			}

			if tc.expectError != "" {
				require.ErrorContains(t, err, tc.expectError, "unexpected error returned")

				return
			}

			require.NoError(t, err, "no error expected")

			assert.Equal(t, tc.expectValue, result, "unexpected value returned")
		})
	}
}
