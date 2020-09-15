#!/usr/bin/env bash

if [ -z "${SERVICE:-}" ]; then
  echo "✋ Missing SERVICE!" >&2
  exit 1
fi

#if [ -z "${TAG:-}" ]; then
#  TAG="$(git describe --dirty)"
#  echo "🎈 Using ${TAG}!" >&2
#fi

TAG=$(git config --get remote.origin.url)
GITHUB_SSH_PREFIX="git@github.com:"
GITHUB_HTTPS_PREFIX="https://github.com"
GIT_SUFFIX=".git"

TAG=${TAG#$GITHUB_SSH_PREFIX}
TAG=${TAG#$GITHUB_HTTPS_PREFIX}
TAG=${TAG%$GIT_SUFFIX}

docker build -t "docker.pkg.github.com/${TAG}/${SERVICE}:latest" --build-arg SERVICE=tgbot .
