package masscan

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

var (
	ErrTempfileExhausted = errors.New("exhausted attempts to allocate a temp file")
)

type Masscan struct {
	cfg Config
}

func (m *Masscan) Run(ctx context.Context) (Report, error) {
	logger := zerolog.Ctx(ctx)

	tmpfile, cleanup, err := tempFile(m.cfg.TempDir, "json")
	if err != nil {
		return Report{}, err
	}

	defer cleanup()

	args := []string{
		"--output-format", "json",
		"--output-filename", tmpfile,
	}

	if m.cfg.MaxRate > 0 {
		args = append(args, "--max-rate", strconv.Itoa(m.cfg.MaxRate))
	}

	if m.cfg.ConfigPath != "" {
		args = append(args, "-c", m.cfg.ConfigPath)
	} else if m.cfg.Config.Configured() {
		config, err := m.cfg.Config.GetValue(ctx)
		if err != nil {
			return Report{}, fmt.Errorf("failed to get masscan config: %w", err)
		}

		conffile, cleanup, err := tempFile(m.cfg.TempDir, "conf")
		if err != nil {
			return Report{}, err
		}

		defer cleanup()

		if err := os.WriteFile(conffile, []byte(config), 0644); err != nil {
			return Report{}, fmt.Errorf("failed to write config: %w", err)
		}

		args = append(args, "-c", conffile)
	}

	report := Report{
		Partial: true,
		MaxRate: m.cfg.MaxRate,
	}

	if m.cfg.Ranges.Configured() {
		ranges, err := m.cfg.Ranges.GetValue(ctx)
		if err != nil {
			return report, fmt.Errorf("failed to get ranges: %w", err)
		}

		report.Ranges = ranges

		args = append(args, ranges...)
	}

	if m.cfg.Ports.Configured() {
		ports, err := m.cfg.Ports.GetValue(ctx)
		if err != nil {
			return report, fmt.Errorf("failed to get ports: %w", err)
		}

		report.Ports = ports

		args = append(args, "-p"+strings.Join(ports, ","))
	}

	logger.Debug().Msgf("prepared command %s %q", m.cfg.BinPath, args)

	var output bytes.Buffer

	cmd := exec.CommandContext(ctx, m.cfg.BinPath, args...)

	cmd.WaitDelay = m.cfg.WaitDelay
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		return report, fmt.Errorf("failed to run command: %w: %s", err, output.String())
	}

	out := output.String()

	logger.Debug().Msgf("command output: %s", out)

	// The last line of output looks like:
	// rate:  0.00-kpps, 100.00% done, waiting -30-secs, found=0
	lastUpdateIndex := strings.LastIndex(out, "found=")
	foundStr := strings.TrimSpace(out[lastUpdateIndex+6:])

	found, err := strconv.Atoi(foundStr)
	if err != nil {
		logger.Debug().Err(err).Msgf("unable to parse command output for found count: '%s', attempting to read report anyways", foundStr)
	} else if found == 0 {
		logger.Debug().Msg("no results found")

		return report, nil
	} else {
		logger.Debug().Msgf("command reports %d ports found", found)
	}

	return m.generateReport(ctx, tmpfile, report)
}

func (m *Masscan) generateReport(ctx context.Context, file string, report Report) (Report, error) {
	logger := zerolog.Ctx(ctx)

	contents, err := os.ReadFile(file)
	if err != nil {
		return report, err
	}

	if err := json.Unmarshal(contents, &report.RawResults); err != nil {
		logger.Debug().Err(err).Str("report_contents", string(contents)).Msg("failed to decode raw report results")

		return report, fmt.Errorf("failed to decode report results: %w", err)
	}

	for _, entry := range report.RawResults {
		for _, port := range entry.Ports {
			if report.Results == nil {
				report.Results = map[string]Results{
					entry.IP: {
						IP:    entry.IP,
						Ports: []Port{port},
					},
				}
			} else {
				result, ok := report.Results[entry.IP]
				if !ok {
					result.IP = entry.IP
				}

				result.Ports = append(result.Ports, port)

				report.Results[entry.IP] = result
			}
		}
	}

	report.Partial = false

	return report, nil
}

func tempFile(dir string, ext string) (string, func(), error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", nil, err
	}

	var err error

	for i := 0; i < 10; i++ {
		var uniq [8]byte

		_, err = rand.Read(uniq[:])
		if err != nil {
			return "", nil, err
		}

		filename := filepath.Join(dir, fmt.Sprintf("masscan-%02x.%s", string(uniq[:]), ext))

		if _, err = os.Stat(filename); err != nil && errors.Is(err, os.ErrNotExist) {
			return filename, func() {
				os.Remove(filename)
			}, nil
		}
	}

	return "", nil, fmt.Errorf("%w: %w", ErrTempfileExhausted, err)
}

func New(_ context.Context, opts ...Option) (*Masscan, error) {
	cfg := newConfig(opts...)

	return &Masscan{
		cfg: cfg,
	}, nil
}
