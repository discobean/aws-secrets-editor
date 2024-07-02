# AWS Secrets Manager Editor

This Go program allows you to fetch a secret from AWS Secrets Manager, 
edit it using your preferred text editor, and update the secret if 
changes are made. You can also specify a version of the secret to edit.

## Features
- Fetch a secret from AWS Secrets Manager
- Edit the secret using the default or specified text editor
- Update the secret if modifications are made
- Specify a version ID or version stage for the secret

## Prerequisites
- Go (1.16+)
- AWS CLI configured with the necessary permissions

## Installation
1. **Clone the repository**:
    ```sh
    git clone https://github.com/discobean/aws-secrets-editor.git
    cd aws-secrets-editor
    ```

2. **Build**:
    ```sh
    make build
    ```

## Usage
### Build the Program
To build the Go program, run:

```sh
./aws-secrets-editor --help
Usage of ./aws-secrets-editor:
  -editor string
    	The editor to use (optional)
  -secretid string
    	The ID or ARN of the secret to edit
  -versionid string
    	The version ID of the secret (optional)
  -versionstage string
    	The version stage of the secret (optional)
```

### Edit a secret
To edit a secret, run the following command:

```sh
export AWS_PROFILE=some-profile
export AWS_REGION=us-east-1
./aws-secrets-editor -secretid my-secret-name
```