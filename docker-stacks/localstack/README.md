# localstack

[LocalStack](https://docs.localstack.cloud/) community edition for emulating AWS services on your machine. Useful for developing against S3, SQS, SNS, DynamoDB, Lambda, etc. without touching the real cloud.

## What's inside

| Service | Image | Default port |
|---|---|---|
| `localstack` | `localstack/localstack:latest` | `4566` (single multiplexed endpoint) |

Persistence is enabled by default — restarting the stack keeps data around. Set `LOCALSTACK_PERSISTENCE=0` in `.env` to disable.

## Start it

```bash
piltover stacks up localstack
```

## Use it

Point the AWS CLI at the local endpoint:

```bash
aws --endpoint-url=http://localhost:4566 s3 mb s3://my-bucket
aws --endpoint-url=http://localhost:4566 s3 ls
```

Or via environment variable:

```bash
export AWS_ENDPOINT_URL=http://localhost:4566
aws s3 ls
```

For SDK usage, set the endpoint in your client config (e.g. `endpointUrl` in JS, `endpoint_url` in Python boto3, `BaseEndpoint` in Go AWS SDK v2).

## Configure services

`LOCALSTACK_SERVICES` is a comma-separated list of services to enable. The default (`s3,sqs,sns,dynamodb,lambda,iam,sts,logs`) covers the most common needs. See [LocalStack services](https://docs.localstack.cloud/user-guide/aws/feature-coverage/) for the full list.

## Reset

```bash
piltover stacks nuke localstack   # wipes persisted state
```
