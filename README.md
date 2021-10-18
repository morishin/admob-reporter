# admob-reporter

## Setup

1. Copy .env.example to .env and fill in the values.

    ```sh
    cp .env.example .env
    ```

1. Download your own OAuth2 credentials from GCP Console: https://console.cloud.google.com/apis/credentials
1. Rename it into `oauth_client_secret.json` and put it next to oauth_client_secret.json.example.

## Run

```sh
$ cd dev
$ go run .
```

## Deploy to AWS

```sh
$ cd aws/cdk
$ cdk deploy
```

