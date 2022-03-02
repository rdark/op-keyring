package main

import (
	"bytes"
	"fmt"
	"github.com/zalando/go-keyring"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func generateSessionToken(opPath string, serviceName string, userName string) (sessionToken string, err error) {
	opSigninArgs := []string{"signin", "my", "--raw"}
	cmd := exec.Command(opPath, opSigninArgs...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	sessionToken = strings.TrimSuffix(stdout.String(), "\n")
	err = keyring.Set(serviceName, userName, sessionToken)
	return sessionToken, err
}

func runOpCmd(opPath string, serviceName string, userName string, cmdArgs []string) error {
	sessionToken, err := keyring.Get(serviceName, userName)
	if err != nil {
		if err.Error() == "secret not found in keyring" {
			sessionToken, err = generateSessionToken(opPath, serviceName, userName)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	opBaseCmdArgs := []string{"--session", sessionToken}
	opCmdArgs := append(opBaseCmdArgs, cmdArgs...)

	cmd := exec.Command(opPath, opCmdArgs...)
	var stderr bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		invalidTokenRegex := regexp.MustCompile("[iI]nvalid session token\n$")
		if invalidTokenRegex.Match(stderr.Bytes()) {
			err = keyring.Delete(serviceName, userName)
			if err != nil {
				return err
			}
			return runOpCmd(opPath, serviceName, userName, cmdArgs)
		}
		fmt.Print(stderr.String())
	}
	return err
}

func main() {

	opPath, err := exec.LookPath("op")
	if err != nil {
		fmt.Println("Could not find `op` binary in $PATH")
		os.Exit(1)
	}

	var serviceName = "op"
	var userName = "session_token"

	args := os.Args[1:]
	err = runOpCmd(opPath, serviceName, userName, args)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		os.Exit(1)
	}
}
