#!/usr/bin/env bash

# Print OpenFlag version.
/app/openflag version

# Run database migrations.
/app/openflag migrate

# Run OpenFlag server.
/app/openflag server
