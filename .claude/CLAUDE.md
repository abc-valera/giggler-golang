# Software Design

There are two approaches for creating a piece of software: code-first (the code is the single source of truth) and spec-first (specs are the single source of truth). The project usually uses a mix of these two.

Each product is a set of artifacts. Artifacts are any _real_ byproducts created during a project's lifecycle, they define the project and live with it.
Artifacts can be short-lived (are temporary), long-lived (stay up-to-date and evolve alongside the project) and derived (generated from the long-lived artifacts).
**Long-lived artifacts are a single source of truth for the project.**

The following artifacts can be used to define a piece of software.

## Code

This a long-lived (sometimes derived) artifact that implements the features and usecases of the application.

All the code is located in the `./src/{assets, cmd, features, shared}` folders:

- `assets` folder contains the files that are imported in from the source code, e.g. logos, fonts, etc.
- `cmd` folder contains the entrypoint/s of the application, e.g. an http web server, an SSR app, a CLI tool, etc. Each entrypoint is a separate binary that should be built and deployed separately.
- `features` folder contains the features of the application which are just a way of grouping the usecases. Each feature is a separate package/module. Each should export types/functions that implement the usecase/s and are consumed by the `cmd` entrypoint/s. Certain features can consume other features, but always in an acyclic way.
- `shared` folder contains the shared code that is used by multiple features but doesn't contain any actual usecases implementations. It can be a db connection for example. Certain shared folders can consume other shared folders, but always in an acyclic way. Shared folders can't consume features in 99.9% of applications. Before implementing a new helper/shared functionality the developer should generally check if it already exists in the shared folder.

The static files served as-is at the site root are located in the `./public` folder.

The build/compilation outputs to the `./build` folder. All the local files, e.g. caches, logs, etc. are located in the `./var` folder. Both these folders are ignored by git.

### On naming

Stuttering type names are generally considered fine - `foo.NewFoo`, `sort.Sorter` or `user.User`, but at the same time this can be named similar to `foo.New`, `sort.Interface` or `user.Model`. The second option is preferred, but depends on the context. The goal is to avoid unnecessary stuttering.

### Code comments style

Doc comments are being rendered into a standalone reference doc (using dedicated tools e.g. godoc, rustdoc, typedoc, etc.) and should be treated like a brief wiki: a human developer or an AI agent which is using, integrating or developing the code will study these derived docs first, and then the source code only if needed to clarify certain things. Generally, the minimal and brief formulations are preferred. Both doc links and internet links are welcomed and should be used when needed.

The following components must have doc comments:

1. The package/module.

   It should provide the information relevant to the package as a whole and generally sets expectations for the package. Linking to other doc comments should be done when needed. For the features packages the doc should describe the business domain of the feature and, if needed, the technical details of the implementation.

   For the shared packages the doc should describe the technical purpose of the package or/and the problem it solves.

2. Usecase type/function.

   Docs should describe the business purpose of the usecase implemented.

   If the usecase type's API is simple, the doc comment can be quite short.

   A usecase function’s doc comment should explain, if not obvious, what the function inputs and/or outputs. Doc comments should not explain internal details such as the algorithm used in the current implementation. Those are best left to comments inside the function body.

Other doc comments are usually optional:

1. Non usecase, unexported or internal types/functions.

   Can describe the purpose of the type/function, but usually can be omitted because the code itself is self-explanatory.

2. Consts/Vars
   A single doc comment can introduce a group of related constants, with individual constants only documented by short end-of-line comments. Sometimes the group needs no doc comment at all. On the other hand, ungrouped constants typically warrant a full doc comment starting with a complete sentence.

   The conventions for variables are the same as those for constants.

The comments should never be linked or mention the short-lived artifacts, past code/artifacts, personnel, or any other ephemeral things. They should only describe the current state of the code.

### Comment Annotations

There are two kinds of annotations that can be left in the code comments:

1. TODOs using the `TODO:` prefix.
2. Some sort of ideas or uncertainties using the `IDEA:` prefix.

## Package manifest

E.g. `go.mod` or `package.json`.

## Task runner

The `./run.sh` file is used as a task runner. It is an alternative for the Makefile.

## git hooks

Vanilla git hooks are used for this. Usually the `pre-commit` and `pre-push` hook are used to lint, check the build, generate derived artifacts, etc.

The custom ones are stored in `.githooks` and should be copied the `.git` folder after cloning.

## Builds

There are two build configurations for any project:

1. dev - enables debugging, fast build, hot reload, etc.
2. release - optimized for the production env

Build scripts, configurations, etc. can be stored either in the `./scripts` or in the `./scripts/build` folder.

## Deployments

There are three deployments types for any project (though some of them can be missing):

1. local - used locally by the developers for development and testing, can use dev or release builds
2. stage - for internal testing, requires release build. Should also be as identical as possible to the prod deployment.
3. prod - an actual deployment for the end users, requires release build.

## Infrastructure configurations

This includes all the tooling required to be able to perform the deployments described above. Thus it is assumed that for each deployment type there will be a separate configurations per each infrastructure tool.

## Application versions

The Semantic Versioning (SemVer) is used. It goes as `release/vMAJOR.MINOR` for branches and `vMAJOR.MINOR.PATCH[RC-<number>]` for tags, where:

- MAJOR - incremented when breaking API changes are introduced
- MINOR - incremented when new features are added
- PATCH - incremented when bug fixes are introduced
- RC (release candidate) - incremented to track the different staging versions for a release candidate

## git repo management

Trunk based development is used for the git repo management. The `main` branch is the trunk where developers merge small, frequent updates. The short-lived feature/<feature-name> branches are created and squashed into main.

The code is deployed automatically via CD when a certain git tag is pushed:

- `v1.0.0-rc1` for the staging deployment
- `v1.0.0` for the production deployment

The logic behind this naming is as follows: when the team decides to create a deployment from a certain commit, it is tagged as a release candidate and deployed to the staging environment first. If the deployment is successful and no bugs are found, then the same commit is tagged as a `v1.0.0` and deployed to the production environment. If some bugs were found then the `v1.0.0-rc2` is created on the commit with fixes.

Meanwhile the development continues in the trunk, everyone pushes to main as often as possible while QAs test the code in stage, or users test it in prod.

For really small projects these tags can be created in the `main` branch, but it's often better to create a separate branch for the release to make it easier to perform bugfixes. The release branch is created from the commit that should be released (`release/v1.0`) and this commit is then tagged as `v1.0.0-rc1` to deploy to stage and as `v1.0.0` to deploy to prod.

If any bugs were found after the release then it should be reproduced and fixed in the `main` branch first, and then cherry-picked to the `release/v1.0` branch and tagged as `v1.0.1-rc1` or immediately `v1.0.1`. If the bug can't be reproduced in the main (e.g. is already fixed by chance during another feature development, or something else happened), then the fix is performed in the release/1.0 branch directly and merged to main if needed. **Note, that this the only case where code is pushed to other branch than trunk.**

## CI/CD configurations

## Secrets

For the locally stored secrets there is a separate folder (`./secrets`) which is ignored by git and describes its secrets files. They differ for different projects, but for all (or almost all) of them the dotenv files are present (`./secrets/prod.env`, `./secrets/stage.env`, `./secrets/local.env`).

## IDE configurations

A `.vscode/settings.json`.

## Additional developer docs

(_TODO: not used for now, to be explored_)

These are the documents/files that are intended to be read by devs or coding AI agents.

Formats:

- DB schema docs with [tbls](https://github.com/k1Low/tbls) that generates docs from an actual db
- SBOM using syft
- https://structurizr.com/

## End user docs

These are the documents/files that are intended to be read by end users or their AI agents.

Formats:

- OpenAPI which can be used as a derived artifact for the code-first or a long-lived for the spec-first, implemented for many programming languages.
- CLI references for the cli applications.
- Changelog.
