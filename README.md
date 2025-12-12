# giggler-golang

_Work in progress_

## Description

<img src="external/designs/svg/logo.svg" width="120" align="left" style="margin-right: 20px">

**Giggler** is a social network based on the idea of sharing and discussing jokes. This repository contains the **REST API** for the social network.

The project utilizes **feature-first** design, where features are first-class citizens. Each feature is a standalone package, that contains all its initialisation, dependecy injection, etc.

## On the architecture

The app uses a **code-first** approach for the documentation, api-spec, database management, etc. This approach was integrated where possible throught the codebase.

TBD

## Local launch

Either `./run.sh run::webapi:dev` or `./run.sh run::webapi:release`

## Development

To update dependencies, run:

```bash
go get -u ./...
go mod tidy
```

## Designs

All the designs are available in [Figma](https://www.figma.com/design/sdu0PTLD3NOxOLNNI1S23f/)
