default:
    just --list ui

setup:
	npm install

build:
	npm run build

# start UI server (MODE="dev|preview")
[arg("MODE", pattern="dev|preview")]
run MODE: build
	#!/usr/bin/env bash
	if [[ "{{MODE}}" == "dev" ]]; then
		npm run dev
	elif [[ "{{MODE}}" == "preview" ]]; then
		npm run preview
	else
		echo "Unknown mode: {{MODE}}"
	fi

lint:
	npm run lint

fmt:
	npx prettier . --check

fmt_fix:
	npx prettier . --write

