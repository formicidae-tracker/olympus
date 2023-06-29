package olympus

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

type generatorCommand struct {
	Filename string `short:"f" long:"filename" description:"filename to export secrets"`
}

type GenerateVAPIDKeysCommand struct {
	generatorCommand
}

type nopCloser struct {
	w io.Writer
}

func (w nopCloser) Close() error { return nil }

func (w nopCloser) Write(data []byte) (int, error) { return w.w.Write(data) }

func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}

func (g generatorCommand) getOutput() (io.WriteCloser, error) {
	if len(g.Filename) > 0 {
		return os.Create(g.Filename)
	}
	return NopCloser(os.Stdout), nil
}

func (c *GenerateVAPIDKeysCommand) Execute([]string) error {
	private, public, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return fmt.Errorf("could not generate VAPID Keys: %w", err)
	}

	output, err := c.getOutput()
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = fmt.Fprintf(output, "OLYMPUS_VAPID_PRIVATE=%s\n", private)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(output, "OLYMPUS_VAPID_PUBLIC=%s\n", public)
	if err != nil {
		return err
	}

	return nil
}

type GenerateSecretsCommand struct {
	generatorCommand
}

func (c *GenerateSecretsCommand) Execute([]string) error {
	secret := make([]byte, 64)
	_, err := rand.Read(secret)
	if err != nil {
		return err
	}

	output, err := c.getOutput()
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = fmt.Fprintf(output, "OLYMPUS_SECRET=%s\n", base64.URLEncoding.EncodeToString(secret))
	if err != nil {
		return err
	}

	return nil
}

func init() {
	parser.AddCommand("generate-vapid-keys",
		"generates a new VAPID key pair.",
		"generates a new VAPID key pair on stdout or file. It should be put in a .env file that will be loaded by the docker image.",
		&GenerateVAPIDKeysCommand{})

	parser.AddCommand("generate-secret",
		"generates a new secret.",
		"generates a new secret on stdout or file. It should be put in a .env file that will be loaded by the docker image. This secret is used for the Nonce generation using HMAC to generate CSRF token.",
		&GenerateSecretsCommand{})

}
