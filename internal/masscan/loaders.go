package masscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func loadValue(value string) string {
	scheme, remain, found := strings.Cut(value, "://")
	if !found {
		return value
	}

	switch scheme {
	case "env":
		return os.Getenv(remain)
	case "file":
		b, _ := os.ReadFile(remain)

		return string(b)
	}

	return value
}

func loadEnv[T any](ctx context.Context, env string) (T, error) {
	env = loadValue(env)

	value := os.Getenv(env)

	return decodeValue[T](ctx, "env", env, []byte(value))
}

func loadFile[T any](ctx context.Context, path string) (T, error) {
	var empty T

	path = loadValue(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return empty, fmt.Errorf("error reading file '%s': %w", path, err)
	}

	return decodeValue[T](ctx, "file", path, data)
}

func loadURL[T any](ctx context.Context, uri string, config URLConfig) (T, error) {
	var empty T

	uri = loadValue(uri)

	url, err := url.Parse(uri)
	if err != nil || url.Host == "" {
		url, err = url.Parse("https://" + uri)
		if err != nil {
			return empty, fmt.Errorf("error parsing url '%s': %w", uri, err)
		}
	}

	var data []byte

	switch url.Scheme {
	case "http", "https":
		req, err := http.NewRequestWithContext(ctx, config.getMethod(), url.String(), config.getBody())
		if err != nil {
			return empty, fmt.Errorf("error creating request for '%s': %w", url.String(), err)
		}

		for k, v := range config.getHeaders() {
			req.Header.Set(k, v)
		}

		username := loadValue(config.Auth.Username)
		password := loadValue(config.Auth.Password)

		if username != "" || password != "" {
			req.SetBasicAuth(username, password)
		}

		bearer := loadValue(config.Auth.Bearer)

		if bearer != "" {
			req.Header.Set("Authorization", "Bearer "+bearer)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return empty, fmt.Errorf("error requesting '%s': %w", url.String(), err)
		}

		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return empty, fmt.Errorf("error reading response for '%s': %w", url.String(), err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return empty, fmt.Errorf("unexpected status code for '%s': status: %d body: %s", url.String(), resp.StatusCode, string(data))
		}
	}

	return decodeValue[T](ctx, "url", url.String(), data)
}

func decodeValue[T any](_ context.Context, kind string, ref string, data []byte) (T, error) {
	var empty T

	data = bytes.TrimSpace(data)

	if len(data) == 0 {
		return empty, nil
	}

	var (
		decoder     func(data []byte, out any) error
		decoderType string

		splitDecoder = func(s string) func(data []byte, out any) error {
			return func(data []byte, out any) error {
				switch v := out.(type) {
				case *[]string:
					if s == "" {
						*v = []string{string(data)}

						return nil
					}

					elems := strings.Split(string(data), s)

					*v = make([]string, 0, len(elems))

					for _, elem := range elems {
						elem = strings.TrimSpace(elem)

						if elem != "" {
							*v = append(*v, elem)
						}
					}
				default:
					return fmt.Errorf("invalid type: %T", out)
				}

				return nil
			}
		}
	)

	switch any(empty).(type) {
	case []string:
		switch data[0] {
		case '[':
			decoder = json.Unmarshal
			decoderType = "json"
		default:
			switch {
			case bytes.ContainsRune(data, ','):
				decoder = splitDecoder(",")
				decoderType = "comma-separated"
			case bytes.ContainsRune(data, '\n'):
				decoder = splitDecoder("\n")
				decoderType = "newline-separated"
			default:
				decoder = splitDecoder("")
				decoderType = "single-value"
			}
		}
	case string:
		switch data[0] {
		case '"':
			decoder = json.Unmarshal
			decoderType = "json"
		default:
			decoder = func(data []byte, out any) error {
				if value, ok := out.(*string); ok {
					*value = string(data)

					return nil
				}

				return fmt.Errorf("invalid type: %T", out)
			}
			decoderType = "raw"
		}
	}

	if decoder == nil {
		return empty, fmt.Errorf("unknown type: %T", empty)
	}

	var ret T

	if err := decoder(data, &ret); err != nil {
		return empty, fmt.Errorf("error decoding %s %s response for '%s': %w", decoderType, kind, ref, err)
	}

	return ret, nil
}
