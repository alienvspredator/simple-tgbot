#!/usr/bin/env sh

if [ -z "${DB_USER}" ]; then
  echo "DB_USER is not provided" >& 2
  exit 1
fi

HOST="${HOST:-"127.0.0.1"}"
PORT="${PORT:-26257}"
DATABASE=${DATABASE:-"defaultdb"}

if [ -n "${PASSWORD}" ]; then
  PASSWORD=":${PASSWORD}"
fi

if [ -n "${OPTIONS}" ]; then
  OPTIONS="?${OPTIONS}"
fi

MIGRATIONS_PATH=$(realpath "$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )/../migrations")

migrate -path="${MIGRATIONS_PATH}" -database "cockroachdb://${DB_USER}${PASSWORD}@${HOST}:${PORT}/${DATABASE}${OPTIONS}" up
