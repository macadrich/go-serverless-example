# go-serverless-example
Serverless example using Golang with DynamoDB

## Prerequisites

- [Node.js & NPM](https://github.com/creationix/nvm)
- [Serverless framework](https://serverless.com/framework/docs/providers/aws/guide/installation/): `npm install -g serverless`
- [Go](https://golang.org/dl/)
- [dep](https://github.com/golang/dep): `brew install dep && brew upgrade dep`

## Quick Start

0. Clone the repo

```
git clone git@github.com:macadrich/go-serverless-example.git
cd go-serverless-example
```

1. Install Go dependencies

```
dep ensure
```

2. Compile functions as individual binaries for deployment package:

```
make
```

> You need to perform this compilation step before deploying.

3. Deploy!

```
serverless deploy
