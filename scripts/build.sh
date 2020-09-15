#!/usr/bin/env bash

if [ -z "${SERVICE:-}" ]; then
  echo "✋ Missing SERVICE!" >&2
  exit 1
fi

if [ -z "${TAG:-}" ]; then
  TAG="$(git describe --dirty)"
  echo "🎈 Using ${TAG}!" >&2
fi

docker build -t simple-tgbot --build-arg SERVICE=tgbot .
