package nak

import (
	"context"
	"io"
	"os"

	"github.com/fiatjaf/cli/v3"
)

var version string = "debug"

var app = &cli.Command{
	Name:                      "nak",
	Suggest:                   true,
	UseShortOptionHandling:    true,
	AllowFlagsAfterArguments:  true,
	Usage:                     "the nostr army knife command-line tool",
	DisableSliceFlagSeparator: true,
	Commands: []*cli.Command{
		req,
		count,
		fetch,
		event,
		decode,
		encode,
		key,
		verify,
		relay,
		bunker,
		serve,
		encrypt,
		decrypt,
	},
	Version: version,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:       "quiet",
			Usage:      "do not print logs and info messages to stderr, use -qq to also not print anything to stdout",
			Aliases:    []string{"q"},
			Persistent: true,
			Action: func(ctx context.Context, c *cli.Command, b bool) error {
				q := c.Count("quiet")
				if q >= 1 {
					log = func(msg string, args ...any) {}
					if q >= 2 {
						stdout = func(_ ...any) (int, error) { return 0, nil }
					}
				}
				return nil
			},
		},
	},
}

func Nak(args []string) ([]byte, error) {
	output, err := captureOutput(func() error {
		return app.Run(context.Background(), args)
	})
	if err != nil {
		stdout(err)
		return output, err
	}

	return output, nil
}

func captureOutput(f func() error) ([]byte, error) {
	// Redirect stdout to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the provided function
	err := f()

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read the captured output from the buffer
	output, _ := io.ReadAll(r)
	return output, err
}
