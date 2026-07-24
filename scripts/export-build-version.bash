#!/usr/bin/env bash

# The script exports the build version for the project.
# The build version can be one of the following:
# - dirty: if the working directory has uncommitted changes
# - tag: if the current commit is tagged
# - commit: if the current commit is not tagged

export-build-version() {
	if [[ -n "$(git status --porcelain)" ]]; then
		export GIGGLER_GOLANG_VERSION="dirty"
		return
	fi

	tagVersion=$(git tag --points-at HEAD | head -n1)
	if [[ -n "$tagVersion" ]]; then
		export GIGGLER_GOLANG_VERSION="${tagVersion}"
		return
	fi

	commitVersion=$(git rev-parse --short=4 HEAD)
	export GIGGLER_GOLANG_VERSION="${commitVersion}"
}
