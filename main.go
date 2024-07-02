package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func getSecret(sess *session.Session, secretID, versionID, versionStage string) (string, error) {
	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	if versionID != "" {
		input.VersionId = aws.String(versionID)
	}

	if versionStage != "" {
		input.VersionStage = aws.String(versionStage)
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}

func updateSecret(sess *session.Session, secretID, secretString string) error {
	svc := secretsmanager.New(sess)
	input := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(secretID),
		SecretString: aws.String(secretString),
	}

	_, err := svc.UpdateSecret(input)
	return err
}

func editFile(filename, editor string) error {
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func fileModified(filename string, originalContent []byte) bool {
	currentContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	return string(currentContent) != string(originalContent)
}

func main() {
	secretID := flag.String("secretid", "", "The ID or ARN of the secret to edit")
	editor := flag.String("editor", "", "The editor to use (optional)")
	versionID := flag.String("versionid", "", "The version ID of the secret (optional)")
	versionStage := flag.String("versionstage", "", "The version stage of the secret (optional)")
	flag.Parse()

	if *secretID == "" {
		log.Fatalf("secretid is required")
	}

	if *editor == "" {
		*editor = os.Getenv("EDITOR")
		if *editor == "" {
			*editor = "vi"
		}
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	secret, err := getSecret(sess, *secretID, *versionID, *versionStage)
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

	tmpFile := filepath.Join(usr.HomeDir, "secret.json")
	err = ioutil.WriteFile(tmpFile, []byte(secret), 0644)
	if err != nil {
		log.Fatalf("Failed to write secret to file: %v", err)
	}
	defer os.Remove(tmpFile)

	originalContent, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	err = editFile(tmpFile, *editor)
	if err != nil {
		log.Fatalf("Failed to open editor: %v", err)
	}

	if fileModified(tmpFile, originalContent) {
		updatedContent, err := ioutil.ReadFile(tmpFile)
		if err != nil {
			log.Fatalf("Failed to read modified file: %v", err)
		}

		err = updateSecret(sess, *secretID, string(updatedContent))
		if err != nil {
			log.Fatalf("Failed to update secret: %v", err)
		}

		fmt.Println("Secret updated successfully.")
	} else {
		fmt.Println("No changes made to the secret.")
	}
}
