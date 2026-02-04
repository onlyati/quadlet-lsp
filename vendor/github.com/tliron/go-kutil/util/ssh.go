package util

import (
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func ExecSSH(host string, port int, username string, key string, stdin io.Reader, command ...string) (string, error) {
	if signer, err := ssh.ParsePrivateKey(StringToBytes(key)); err == nil {
		config := ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}

		if client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &config); err == nil {
			defer client.Close()

			if session, err := client.NewSession(); err == nil {
				defer session.Close()

				var stdout strings.Builder
				var stderr strings.Builder

				session.Stdin = stdin
				session.Stdout = &stdout
				session.Stderr = &stderr

				// TODO: handle spaces
				if err := session.Run(strings.Join(command, " ")); err == nil {
					return stdout.String(), nil
				} else {
					// TODO: special error object to contain stderr separately
					return "", fmt.Errorf("%s\n%s", err.Error(), stderr.String())
				}
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func CopySSH(host string, port int, username string, key string, reader io.Reader, targetPath string, permissions *int64) error {
	dir := filepath.Dir(targetPath)
	if _, err := ExecSSH(host, port, username, key, nil, "mkdir", "--parents", dir); err == nil {
		if _, err := ExecSSH(host, port, username, key, reader, "cp", "/dev/stdin", targetPath); err == nil {
			if permissions != nil {
				_, err := ExecSSH(host, port, username, key, nil, "chmod", strconv.FormatInt(*permissions, 8), targetPath)
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
