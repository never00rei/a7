package crypto

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
	"filippo.io/age/agessh"
	"golang.org/x/crypto/ssh"
)

var errMissingSSHKey = errors.New("ssh key path is required for encryption")

func MaybeEncryptBody(body string, enabled bool, sshKeyPath string) (string, bool, error) {
	if !enabled {
		return body, false, nil
	}
	encrypted, err := EncryptBody(body, sshKeyPath)
	if err != nil {
		return "", true, err
	}
	return encrypted, true, nil
}

func EncryptBody(body, sshKeyPath string) (string, error) {
	if sshKeyPath == "" {
		return "", errMissingSSHKey
	}
	recipient, err := recipientFromKeyFile(sshKeyPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	armorWriter := armor.NewWriter(&buf)
	enc, err := age.Encrypt(armorWriter, recipient)
	if err != nil {
		return "", fmt.Errorf("encrypt body: %w", err)
	}
	if _, err := io.WriteString(enc, body); err != nil {
		_ = enc.Close()
		_ = armorWriter.Close()
		return "", fmt.Errorf("encrypt body: %w", err)
	}
	if err := enc.Close(); err != nil {
		_ = armorWriter.Close()
		return "", fmt.Errorf("encrypt body: %w", err)
	}
	if err := armorWriter.Close(); err != nil {
		return "", fmt.Errorf("encrypt body: %w", err)
	}
	return buf.String(), nil
}

func DecryptBody(body, sshKeyPath string) (string, error) {
	if sshKeyPath == "" {
		return "", errMissingSSHKey
	}
	identity, err := identityFromKeyFile(sshKeyPath)
	if err != nil {
		return "", err
	}
	reader := strings.NewReader(body)
	armorReader := armor.NewReader(reader)
	dec, err := age.Decrypt(armorReader, identity)
	if err != nil {
		return "", fmt.Errorf("decrypt body: %w", err)
	}
	plain, err := io.ReadAll(dec)
	if err != nil {
		return "", fmt.Errorf("decrypt body: %w", err)
	}
	return string(plain), nil
}

func recipientFromKeyFile(path string) (age.Recipient, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read ssh key: %w", err)
	}
	if recipient, err := agessh.ParseRecipient(strings.TrimSpace(string(data))); err == nil {
		return recipient, nil
	}
	signer, err := ssh.ParsePrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("parse ssh key: %w", err)
	}
	publicKey := signer.PublicKey()
	switch publicKey.Type() {
	case "ssh-ed25519":
		recipient, err := agessh.NewEd25519Recipient(publicKey)
		if err != nil {
			return nil, fmt.Errorf("parse ssh key: %w", err)
		}
		return recipient, nil
	case "ssh-rsa":
		recipient, err := agessh.NewRSARecipient(publicKey)
		if err != nil {
			return nil, fmt.Errorf("parse ssh key: %w", err)
		}
		return recipient, nil
	default:
		return nil, fmt.Errorf("unsupported ssh key type: %s", publicKey.Type())
	}
}

func identityFromKeyFile(path string) (age.Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read ssh key: %w", err)
	}
	identity, err := agessh.ParseIdentity(data)
	if err != nil {
		return nil, fmt.Errorf("parse ssh key: %w", err)
	}
	return identity, nil
}
