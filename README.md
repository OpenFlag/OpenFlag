# OpenFlag [Under Development]

[![Build Status][1]][2]
[![CodeCov][3]][4]
[![MIT Licence][5]][6]
[![Go Report][7]][8]
[![PkgGoDev][9]][10]

# Introduction

OpenFlag is an open-source feature flagging, A/B testing, and dynamic configuration service. It delivers the right experience to the right entity and monitors the impact. It has clear swagger REST APIs for flags management and flag evaluation.

## Documentation

You can find documentation in <a href="https://openflag.github.io">here</a>.

## Quick demo

Try it with Docker.

```bash
# Download docker-compose file.
wget https://raw.githubusercontent.com/OpenFlag/OpenFlag/master/docker-compose.yml

# Start using docker-compose.
docker-compose up -d

# Open the OpenFlag UI and create your feature flag, experiment, or configuration.
open 127.0.0.1:7677

# Sending a request for evaluation in the UI panel or using curl.
curl --location --request POST 'http://127.0.0.1:7677/api/v1/evaluation' \
--header 'Content-Type: application/json' \
--data-raw '{
    "entities": [
        {
            "entity_id": 1234567,
            "entity_type": "type1",
            "entity_context": {
                "state": "CA"
            }
        }
    ],
    "flags": ["flag1", "flag2"],
    "save_contexts": true,
    "use_stored_contexts": false
}'
```

### TODO

* Support custom constraints using **pluggable** languages.

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
