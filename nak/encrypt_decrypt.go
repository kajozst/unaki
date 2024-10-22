package nak

import (
	"context"
	"fmt"

	"github.com/fiatjaf/cli/v3"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
)

var encrypt = &cli.Command{
	Name:                      "encrypt",
	Usage:                     "encrypts a string with nip44 (or nip04 if specified using a flag) and returns the resulting ciphertext as base64",
	ArgsUsage:                 "[plaintext string]",
	DisableSliceFlagSeparator: true,
	Flags: append(
		defaultKeyFlags,
		&cli.StringFlag{
			Name:     "recipient-pubkey",
			Aliases:  []string{"p", "tgt", "target", "pubkey"},
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "nip04",
			Usage: "use nip04 encryption instead of nip44",
		},
	),
	Action: func(ctx context.Context, c *cli.Command) error {
		target := c.String("recipient-pubkey")
		if !nostr.IsValidPublicKey(target) {
			return fmt.Errorf("target %s is not a valid public key", target)
		}

		plaintext := c.Args().First()

		if c.Bool("nip04") {
			sec, bunker, err := gatherSecretKeyOrBunkerFromArguments(ctx, c)
			if err != nil {
				return err
			}

			if bunker != nil {
				ciphertext, err := bunker.NIP04Encrypt(ctx, target, plaintext)
				if err != nil {
					return err
				}
				stdout(ciphertext)
			} else {
				ss, err := nip04.ComputeSharedSecret(target, sec)
				if err != nil {
					return fmt.Errorf("failed to compute nip04 shared secret: %w", err)
				}
				ciphertext, err := nip04.Encrypt(plaintext, ss)
				if err != nil {
					return fmt.Errorf("failed to encrypt as nip04: %w", err)
				}
				stdout(ciphertext)
			}
		} else {
			kr, err := gatherKeyerFromArguments(ctx, c)
			if err != nil {
				return err
			}

			res, err := kr.Encrypt(ctx, plaintext, target)
			if err != nil {
				return fmt.Errorf("failed to encrypt: %w", err)
			}
			stdout(res)
		}

		return nil
	},
}

var decrypt = &cli.Command{
	Name:                      "decrypt",
	Usage:                     "decrypts a base64 nip44 ciphertext (or nip04 if specified using a flag) and returns the resulting plaintext",
	ArgsUsage:                 "[ciphertext base64]",
	DisableSliceFlagSeparator: true,
	Flags: append(
		defaultKeyFlags,
		&cli.StringFlag{
			Name:     "sender-pubkey",
			Aliases:  []string{"p", "src", "source", "pubkey"},
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "nip04",
			Usage: "use nip04 encryption instead of nip44",
		},
	),
	Action: func(ctx context.Context, c *cli.Command) error {
		source := c.String("sender-pubkey")
		if !nostr.IsValidPublicKey(source) {
			return fmt.Errorf("source %s is not a valid public key", source)
		}

		ciphertext := c.Args().First()

		if c.Bool("nip04") {
			sec, bunker, err := gatherSecretKeyOrBunkerFromArguments(ctx, c)
			if err != nil {
				return err
			}

			if bunker != nil {
				plaintext, err := bunker.NIP04Decrypt(ctx, source, ciphertext)
				if err != nil {
					return err
				}
				stdout(plaintext)
			} else {
				ss, err := nip04.ComputeSharedSecret(source, sec)
				if err != nil {
					return fmt.Errorf("failed to compute nip04 shared secret: %w", err)
				}
				plaintext, err := nip04.Decrypt(ciphertext, ss)
				if err != nil {
					return fmt.Errorf("failed to encrypt as nip04: %w", err)
				}
				stdout(plaintext)
			}
		} else {
			kr, err := gatherKeyerFromArguments(ctx, c)
			if err != nil {
				return err
			}

			res, err := kr.Decrypt(ctx, ciphertext, source)
			if err != nil {
				return fmt.Errorf("failed to encrypt: %w", err)
			}
			stdout(res)
		}

		return nil
	},
}
