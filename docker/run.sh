#!/usr/bin/env bash

/app/openflag migrate --path=/app/migrations
/app/openflag server
