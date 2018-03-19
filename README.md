# Go Chi Framework

## Status

[Changelog](./CHANGELOG.md)

Build | Status
--- | ---
master | [![Build Status](https://circleci.com/gh/jsm/go-chi-framework/tree/master.svg?style=shield&circle-token=50403245ada8fbacb76e7e50487a9f89afb4beb1)](https://circleci.com/gh/jsm/go-chi-framework/tree/master)
deploy/api | [![Build Status](https://circleci.com/gh/jsm/go-chi-framework/tree/deploy%2Fapi.svg?style=shield&circle-token=50403245ada8fbacb76e7e50487a9f89afb4beb1)](https://circleci.com/gh/jsm/go-chi-framework/tree/deploy%2Fapi)
deploy/api-hotfix | [![Build Status](https://circleci.com/gh/jsm/go-chi-framework/tree/deploy%2Fapi-hotfix.svg?style=shield&circle-token=50403245ada8fbacb76e7e50487a9f89afb4beb1)](https://circleci.com/gh/jsm/go-chi-framework/tree/deploy%2Fapi-hotfix)

## About

This is a framework developed by myself, @rohenp, and @Thorbenandresen while at @letshuddleup

Built on top of the following technologies

- https://github.com/go-chi/chi (lightweight go router)
- https://github.com/RichardKnop/machinery (asynchronous task queue/job queue)
- Alembic - http://alembic.zzzcomputing.com/en/latest/ (SQL Migration Tool)
- PostgreSQL (for relational storage)
- Redis (for task queuing)
- AWS Elastic Beanstalk (for deployment)

## Setup
Setup Environment Variables

```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Clone to your Go src

```
mkdir -p $GOPATH/src/github.com/jsm
cd $GOPATH/src/github.com/jsm
git clone git@github.com:jsm/go-chi-framework.git gode
```

Install Go

```
brew install go
```

Run Setup

```
cd $GOPATH/src/github.com/jsm/gode
make setup
```

### Local Development

#### Running locally

##### Redis

First spin up supporting software

```
make run-support
```

##### Worker

```
make run-worker
```

##### API

Finally, run the API

```
make run-api
```

#### Running in docker

Alternatively, you can run all services at once with docker with

```
make docker
```

### Deployment to Dev

In one commit to master:

1. Update VERSION file. Follow [Semantic Versioning](http://semver.org/)
2. Ensure "Next" in [Changelog](./CHANGELOG.md) is up to date
3. Change "Next" to the version being deployed
4. Add a new "Next" Section. You can copy [this](./NEXT.md)
5. Commit your changes with the message `Version X.X.X`
6. Push to master

Whenever the master build runs with a new VERSION, it will deploy it to api-dev

### Deploying a hotfix

```
make deploy-api-hotfix
```

#### Promotion to Prod, Switching Version, Or Rollback

Versions are tagged by their VERSION numbers

- Go to elastic beanstalk
- Click the version you want to deploy
- Hit the deploy button
- Select the environment you want to deploy that version to

## Running Tests

```
make run-tests
```
