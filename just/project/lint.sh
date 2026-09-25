#!/usr/bin/env bash
set -euo pipefail

gocyclo -top 20 -avg -ignore '_test\.go$|\.gen\.go$|/\.build/' .
ineffassign ./...
