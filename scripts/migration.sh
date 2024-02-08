#!/bin/sh

export GOOSE_DRIVER=sqlite3
export GOOSE_DBSTRING=./data.db

case "$1" in
    make) goose -dir migrations create "$2" sql ;;
    status) goose -dir migrations status ;;
    down) goose -dir migrations down ;;
    up) goose -dir migrations up ;;
esac
