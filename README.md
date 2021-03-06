# OpenFlag

[![Build Status][1]][2]
[![Code Cov][3]][4]
[![MIT Licence][5]][6]
[![Go Report][7]][8]
[![Pkg Go Dev][9]][10]
[![Docker Pulls][11]][12]

# Introduction

OpenFlag is an open-source feature flagging, A/B testing, and dynamic configuration service. It delivers the right experience to the right entity and monitors the impact. It has clear Swagger REST APIs for flag management and evaluation.

## Features

* Cloud-native and Kubernetes compatible.
* High performance and easily scalable.
* Support gRPC for flag evaluation.
* Clear Swagger REST APIs for flag management and flag evaluation.
* Rule engine and user segmentation using algebra expression as simple as possible for defining complicated flags.
* Showing the history of a flag.
* Evaluation logging for your data pipeline.
* Contexts saving and reuse stored contexts.
* Support Feature Flagging, Experimentation A/B testing, and Dynamic Configuration.

## Documentation

You can find documentation in <a href="https://openflag.github.io">here</a>.

## Quick demo

Try it with Docker.

```bash
# Download docker-compose file.
wget https://raw.githubusercontent.com/OpenFlag/OpenFlag/master/deployments/docker/openflag/docker-compose.yml

# Start using docker-compose.
docker-compose up -d

# Open the OpenFlag UI and create your feature flag, experiment, or configuration.
open 127.0.0.1:7677
# Note: It doesn't work now because the UI panel is under development. You can use the Swagger file for working with APIs instead.

# Sending a request for evaluation in the UI panel or using curl.
curl --location --request POST 'http://127.0.0.1:7677/api/v1/evaluation' \
--header 'Content-Type: application/json' \
--data-raw '{
    "entities": [
        {
            "entity_id": 1234567,
            "entity_type": "user",
            "entity_context": {
                "city": "AMS"
            }
        }
    ],
    "flags": ["flag1", "flag2"],
}'
```

[1]: https://img.shields.io/drone/build/OpenFlag/OpenFlag.svg?style=flat-square&logo=drone
[2]: https://cloud.drone.io/OpenFlag/OpenFlag
[3]: https://img.shields.io/codecov/c/gh/OpenFlag/OpenFlag?logo=codecov&style=flat-square
[4]: https://codecov.io/gh/OpenFlag/OpenFlag
[5]: https://img.shields.io/github/license/OpenFlag/OpenFlag?style=flat-square
[6]: https://opensource.org/licenses/mit-license.php
[7]: https://goreportcard.com/badge/github.com/OpenFlag/OpenFlag?style=flat-square
[8]: https://goreportcard.com/report/github.com/OpenFlag/OpenFlag
[9]: https://pkg.go.dev/badge/github.com/OpenFlag/OpenFlag
[10]: https://pkg.go.dev/github.com/OpenFlag/OpenFlag
[11]: https://img.shields.io/docker/pulls/openflag/openflag.svg?style=flat-square
[12]: https://hub.docker.com/r/openflag/openflag
