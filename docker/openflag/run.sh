#!/usr/bin/env bash

# Print OpenFlag version.
/app/openflag version

# Run database migrations.
/app/openflag migrate
if [[ $? -ne 0 ]] ; then
    exit 1
fi

# Run OpenFlag server.
/app/openflag server
