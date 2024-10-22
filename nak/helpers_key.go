package nak

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/fiatjaf/cli/v3"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/keyer"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/nbd-wtf/go-nostr/nip46"
	"github.com/nbd-wtf/go-nostr/nip49"
)

var defaultKeyFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "sec",
		Usage:       "secret key to sign the event, as nsec, ncryptsec or hex, or a bunker URL",
		DefaultText: "the key '1'",
		Aliases:     []string{"connect"},
		Category:    CATEGORY_SIGNER,
	},
	&cli.BoolFlag{
		Name:     "prompt-sec",
		Usage:    "prompt the user to paste a hex or nsec with which to sign the event",
		Category: CATEGORY_SIGNER,
	},
	&cli.StringFlag{
		Name:        "connect-as",
		Usage:       "private key to when communicating with the bunker given on --connect",
		DefaultText: "a random key",
		Category:    CATEGORY_SIGNER,
	},
}

func gatherKeyerFromArguments(ctx context.Context, c *cli.Command) (keyer.Keyer, error) {
	key, bunker, err := gatherSecretKeyOrBunkerFromArguments(ctx, c)
	if err != nil {
		return nil, err
	}

	var kr keyer.Keyer
	if bunker != nil {
		kr = keyer.NewBunkerSignerFromBunkerClient(bunker)
	} else {
		kr, err = keyer.NewPlainKeySigner(key)
	}

	return kr, err
}

func gatherSecretKeyOrBunkerFromArguments(ctx context.Context, c *cli.Command) (string, *nip46.BunkerClient, error) {
	var err error

	sec := c.String("sec")
	if strings.HasPrefix(sec, "bunker://") {
		// it's a bunker
		bunkerURL := sec
		clientKey := c.String("connect-as")
		if clientKey != "" {
			clientKey = strings.Repeat("0", 64-len(clientKey)) + clientKey
		} else {
			clientKey = nostr.GeneratePrivateKey()
		}
		bunker, err := nip46.ConnectBunker(ctx, clientKey, bunkerURL, nil, func(s string) {
			log(color.CyanString("[nip46]: open the following URL: %s"), s)
		})
		return "", bunker, err
	}

	// take private from flags, environment variable or default to 1
	if sec == "" {
		if key, ok := os.LookupEnv("NOSTR_SECRET_KEY"); ok {
			sec = key
		} else {
			sec = "0000000000000000000000000000000000000000000000000000000000000001"
		}
	}

	if c.Bool("prompt-sec") {
		if isPiped() {
			return "", nil, fmt.Errorf("can't prompt for a secret key when processing data from a pipe, try again without --prompt-sec")
		}
		sec, err = askPassword("type your secret key as ncryptsec, nsec or hex: ", nil)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get secret key: %w", err)
		}
	}

	if strings.HasPrefix(sec, "ncryptsec1") {
		sec, err = promptDecrypt(sec)
		if err != nil {
			return "", nil, fmt.Errorf("failed to decrypt: %w", err)
		}
	} else if bsec, err := hex.DecodeString(leftPadKey(sec)); err == nil {
		sec = hex.EncodeToString(bsec)
	} else if prefix, hexvalue, err := nip19.Decode(sec); err != nil {
		return "", nil, fmt.Errorf("invalid nsec: %w", err)
	} else if prefix == "nsec" {
		sec = hexvalue.(string)
	}

	if ok := nostr.IsValid32ByteHex(sec); !ok {
		return "", nil, fmt.Errorf("invalid secret key")
	}

	return sec, nil, nil
}

func promptDecrypt(ncryptsec string) (string, error) {
	for i := 1; i < 4; i++ {
		var attemptStr string
		if i > 1 {
			attemptStr = fmt.Sprintf(" [%d/3]", i)
		}
		password, err := askPassword("type the password to decrypt your secret key"+attemptStr+": ", nil)
		if err != nil {
			return "", err
		}
		sec, err := nip49.Decrypt(ncryptsec, password)
		if err != nil {
			continue
		}
		return sec, nil
	}
	return "", fmt.Errorf("couldn't decrypt private key")
}

func askPassword(msg string, shouldAskAgain func(answer string) bool) (string, error) {
	config := &readline.Config{
		Stdout:                 color.Error,
		Prompt:                 color.YellowString(msg),
		InterruptPrompt:        "^C",
		DisableAutoSaveHistory: true,
		EnableMask:             true,
		MaskRune:               '*',
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return "", err
	}

	for {
		answer, err := rl.Readline()
		if err != nil {
			return "", err
		}
		answer = strings.TrimSpace(answer)
		if shouldAskAgain != nil && shouldAskAgain(answer) {
			continue
		}
		return answer, err
	}
}
