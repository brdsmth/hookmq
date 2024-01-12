## Introduction

**HookMQ** emerged to efficiently manage job queues, including those with future execution dates.

## Architecture

## Example

```javascript
POST https://localhost:8081/api/queue HTTP/1.1
content-type: application/json

{
    payload: { /* customize the payload how you want */ },
    url: "https://your.server.com/some-action",
    executeAt: "2024-01-10T02:36:00Z"
}

```

## Local Development

#### Environment

The following environment variables need to be set in order to run `hookmq` locally.

```
PORT=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_DEFAULT_REGION=
DYNAMODB_QUEUE_TABLE=
DYNAMODB_PROCESSED_TABLE=
SQS_URL=
SQS_DL_URL=
```

### Usage

```
go run main.go
```

## Deployment

### Docker

#### Build

```
docker build -t hookmq:latest .
```

#### Run

```
docker rm -f hookmq
docker run -p 8081:8081 \
-v ~/.aws/config:/root/.aws/config \
-v ~/.aws/credentials:/root/.aws/credentials \
--name hookmq hookmq:latest
```

The first command force removes the existing docker container and the second command starts a new one. `-p 8081:8081` tells `docker` to map port 8081 on the host to port 8081 inside the container. The connections to AWS resources requires active AWS credentials and the `docker` command includes instructions to mount the local AWS configuration to the docker build.

### Heroku

This section is to help with deploying `hookmq` to heroku.

#### Create application in Heroku

```
heroku create hookmq
```

When deploying to Heroku, we don't have the same level of control over the file system as when running a container locally with Docker. We can't directly replicate the volume mounting (-v) part of the Docker command because Heroku's container runtime doesn't allow us to mount host directories into our dyno.

We can achieve a similar result by setting AWS configuration through environment (config) variables in Heroku.

```
heroku config:set AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -a hookmq
heroku config:set AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY -a hookmq
heroku config:set AWS_DEFAULT_REGION=$AWS_DEFAULT_REGION -a hookmq
```

#### Build and push

```
heroku container:push web -a hookmq
```

#### Release

```
heroku container:release web -a hookmq
```

#### Open

```
heroku open -a hookmq
```

#### Logs

```
heroku logs --tail -a hookmq
```

## Inspiration

This project was mainly coded to following albums:
