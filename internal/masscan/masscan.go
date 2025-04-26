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
	logger zerolog.Logger
	cfg    Config
}

func (m *Masscan) Run(ctx context.Context) (Report, error) {
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
	} else if m.cfg.Config != "" {
		conffile, cleanup, err := tempFile(m.cfg.TempDir, "conf")
		if err != nil {
			return Report{}, err
		}

		defer cleanup()

		if err := os.WriteFile(conffile, []byte(m.cfg.Config), 0x644); err != nil {
			return Report{}, fmt.Errorf("failed to write config: %w", err)
		}

		args = append(args, "-c", conffile)
	}

	args = append(args, m.cfg.Ranges...)

	if len(m.cfg.Ports) != 0 {
		ports := strings.Join(m.cfg.Ports, ",")

		args = append(args, "-p"+ports)
	}

	m.logger.Debug().Msgf("prepared command %s %q", m.cfg.BinPath, args)

	var output bytes.Buffer

	cmd := exec.CommandContext(ctx, m.cfg.BinPath, args...)

	cmd.WaitDelay = m.cfg.WaitDelay
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		return Report{}, fmt.Errorf("failed to run command: %w: %s", err, output.String())
	}

	return m.generateReport(tmpfile)
}

func (m *Masscan) generateReport(file string) (Report, error) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return Report{}, err
	}

	report := Report{
		Ranges:  m.cfg.Ranges,
		Ports:   m.cfg.Ports,
		MaxRate: m.cfg.MaxRate,
	}

	if err := json.Unmarshal(contents, &report.RawResults); err != nil {
		return Report{}, err
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

	return report, nil
}

func tempFile(dir string, ext string) (string, func(), error) {
	if err := os.MkdirAll(dir, 0x755); err != nil {
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

func New(ctx context.Context, opts ...Option) (*Masscan, error) {
	cfg := newConfig(opts...)

	logger := zerolog.Ctx(ctx).With().Str("component", "masscan").Logger()

	return &Masscan{
		logger: logger,
		cfg:    cfg,
	}, nil
}
