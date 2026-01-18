package masscan

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var reValidEnv = regexp.MustCompile(`[^A-Z0-9_]+`)

func testEnvValue(t *testing.T, envSuffix, value string) string {
	t.Helper()

	envName := strings.ToUpper(t.Name())

	if envSuffix != "" {
		envName += "_" + strings.ToUpper(envSuffix)
	}

	env := reValidEnv.ReplaceAllString(envName, "_")

	os.Setenv(env, value)

	t.Cleanup(func() {
		os.Unsetenv(env)
	})

	return env
}

func testFileValue(t *testing.T, value string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.file")

	err := os.WriteFile(path, []byte(value), 0644)

	require.NoError(t, err, "no error expected writing test file")

	return path
}

func TestDynamicValue_Configured(t *testing.T) {
	t.Parallel()

	t.Run("string value", func(t *testing.T) {
		t.Parallel()

		var value DynamicValue[string]

		assert.False(t, value.Configured(), "expected empty struct to be not configured")

		value.Value = "something"

		assert.True(t, value.Configured(), "expected populated struct to be configured")
	})

	t.Run("string slice value", func(t *testing.T) {
		t.Parallel()

		var value DynamicValue[[]string]

		assert.False(t, value.Configured(), "expected empty struct to be not configured")

		value.Value = []string{"something"}

		assert.True(t, value.Configured(), "expected populated struct to be configured")

		value.Value = []string{}

		assert.False(t, value.Configured(), "expected slice length of 0 to be not configured")
	})

	t.Run("env", func(t *testing.T) {
		t.Parallel()

		value := DynamicValue[string]{
			Env: "SOME_ENV",
		}

		assert.True(t, value.Configured(), "expected populated struct to be configured")
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()

		value := DynamicValue[string]{
			File: "some/file",
		}

		assert.True(t, value.Configured(), "expected populated struct to be configured")
	})

	t.Run("url", func(t *testing.T) {
		t.Parallel()

		value := DynamicValue[string]{
			URL: "http://some.example.com",
		}

		assert.True(t, value.Configured(), "expected populated struct to be configured")
	})

	t.Run("url config", func(t *testing.T) {
		t.Parallel()

		value := DynamicValue[string]{
			URLConfig: URLConfig{
				Method: http.MethodPost,
			},
		}

		assert.False(t, value.Configured(), "expected url config to not contribute to configured status")
	})
}

func TestDynamicValue_GetValue(t *testing.T) {
	t.Parallel()

	httpResp := map[string]string{}

	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, ok := httpResp[r.URL.String()]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))

			return
		}

		w.Write([]byte(resp))
	}))

	t.Cleanup(httpSrv.Close)

	testHTTPValue := func(t *testing.T, subPath, retBody string) string {
		t.Helper()

		path := "/" + t.Name() + "/" + subPath

		httpResp[path] = retBody

		return httpSrv.URL + path
	}

	testCases := []struct {
		name         string
		stringConfig *DynamicValue[string]
		sliceConfig  *DynamicValue[[]string]
		expectValue  any
		expectError  string
	}{
		{
			name:         "no config returns empty string",
			stringConfig: &DynamicValue[string]{},
			expectValue:  "",
		},
		{
			name:        "no config returns nil slice",
			sliceConfig: &DynamicValue[[]string]{},
			expectValue: []string(nil),
		},
		{
			name: "empty slice returns nil slice",
			sliceConfig: &DynamicValue[[]string]{
				Value: []string{},
			},
			expectValue: []string(nil),
		},
		{
			name: "static string",
			stringConfig: &DynamicValue[string]{
				Value: "some value",
			},
			expectValue: "some value",
		},
		{
			name: "static slice",
			sliceConfig: &DynamicValue[[]string]{
				Value: []string{"some value"},
			},
			expectValue: []string{"some value"},
		},
		{
			name: "env string",
			stringConfig: &DynamicValue[string]{
				Env: testEnvValue(t, "env string", "env string"),
			},
			expectValue: "env string",
		},
		{
			name: "env slice",
			sliceConfig: &DynamicValue[[]string]{
				Env: testEnvValue(t, "env slice", "env,slice"),
			},
			expectValue: []string{"env", "slice"},
		},
		{
			name: "file string",
			stringConfig: &DynamicValue[string]{
				File: testFileValue(t, "file string"),
			},
			expectValue: "file string",
		},
		{
			name: "file slice",
			sliceConfig: &DynamicValue[[]string]{
				File: testFileValue(t, "env,slice"),
			},
			expectValue: []string{"env", "slice"},
		},
		{
			name: "http string",
			stringConfig: &DynamicValue[string]{
				URL: testHTTPValue(t, "http-string", "http string"),
			},
			expectValue: "http string",
		},
		{
			name: "http slice",
			sliceConfig: &DynamicValue[[]string]{
				URL: testHTTPValue(t, "http-slice", "env,slice"),
			},
			expectValue: []string{"env", "slice"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				value any
				err   error
			)

			switch {
			case tc.stringConfig != nil && tc.sliceConfig != nil:
				assert.FailNow(t, "stringConfig and sliceConfig in test are mutually exclusive")
			case tc.stringConfig != nil:
				value, err = tc.stringConfig.GetValue(t.Context())
			case tc.sliceConfig != nil:
				value, err = tc.sliceConfig.GetValue(t.Context())
			default:
				assert.FailNow(t, "test case has no stringConfig or sliceConfig set")
			}

			if tc.expectError != "" {
				require.ErrorContains(t, err, tc.expectError, "unexpected error returned from GetValue")

				return
			}

			require.NoError(t, err, "no error expected to be returned from GetValue")

			assert.Equal(t, tc.expectValue, value, "unexpected value returned from GetValue")
		})
	}
}

func TestDynamicValue_UnmarshalMapstructure(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       any
		expect      any
		expectError string
	}{
		{
			"static string",
			"some value",
			DynamicValue[string]{
				Value: "some value",
			},
			"",
		},
		{
			"static slice",
			[]any{"some", "value"},
			DynamicValue[[]string]{
				Value: []string{"some", "value"},
			},
			"",
		},
		{
			"invalid static value",
			123,
			DynamicValue[string]{},
			"expected type 'string', got unconvertible type 'int'",
		},
		{
			"unknown scheme static value",
			"my://static-value",
			DynamicValue[string]{
				Value: "my://static-value",
			},
			"",
		},
		{
			"env string (dynamic string)",
			"env://SOME_ENV",
			DynamicValue[string]{
				Env: "SOME_ENV",
			},
			"",
		},
		{
			"file string (dynamic string)",
			"file://some/file",
			DynamicValue[string]{
				File: "some/file",
			},
			"",
		},
		{
			"http string (dynamic string)",
			"http://some.example.com",
			DynamicValue[string]{
				URL: "http://some.example.com",
			},
			"",
		},
		{
			"https string (dynamic string)",
			"https://some.example.com",
			DynamicValue[string]{
				URL: "https://some.example.com",
			},
			"",
		},
		{
			"env string (dynamic struct)",
			map[string]any{
				"env": "SOME_ENV",
			},
			DynamicValue[string]{
				Env: "SOME_ENV",
			},
			"",
		},
		{
			"file string (dynamic struct)",
			map[string]any{
				"file": "some/file",
			},
			DynamicValue[string]{
				File: "some/file",
			},
			"",
		},
		{
			"http string (dynamic struct)",
			map[string]any{
				"url": "http://some.example.com",
			},
			DynamicValue[string]{
				URL: "http://some.example.com",
			},
			"",
		},
		{
			"https string (dynamic struct)",
			map[string]any{
				"url": "https://some.example.com",
			},
			DynamicValue[string]{
				URL: "https://some.example.com",
			},
			"",
		},
		{
			"url with config (dynamic struct)",
			map[string]any{
				"url": "https://some.example.com",
				"url_config": map[string]any{
					"method": "POST",
				},
			},
			DynamicValue[string]{
				URL: "https://some.example.com",
				URLConfig: URLConfig{
					Method: "POST",
				},
			},
			"",
		},
		{
			"struct value type",
			map[string]any{
				"key": "some value",
			},
			DynamicValue[struct{ Key string }]{
				Value: struct{ Key string }{
					Key: "some value",
				},
			},
			"",
		},
		{
			"not dynamic struct",
			map[string]any{
				"key":  "some value",
				"file": "some file",
			},
			DynamicValue[struct{ Key string }]{},
			"has invalid keys",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				result any

				// since result is a pointer, we'll get the pointer to the expect value.
				expectPtr any
			)

			switch expect := tc.expect.(type) {
			case DynamicValue[string]:
				result = &DynamicValue[string]{}
				expectPtr = &expect
			case DynamicValue[[]string]:
				result = &DynamicValue[[]string]{}
				expectPtr = &expect
			case DynamicValue[struct{ Key string }]:
				result = &DynamicValue[struct{ Key string }]{}
				expectPtr = &expect
			}

			err := result.(mapstructure.Unmarshaler).UnmarshalMapstructure(tc.input)

			if tc.expectError != "" {
				require.ErrorContains(t, err, tc.expectError, "unexpected error returned")

				return
			}

			require.NoError(t, err, "no error expected")

			assert.Equal(t, expectPtr, result, "unexpected result")
		})
	}
}
