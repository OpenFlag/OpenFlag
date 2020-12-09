#!/usr/bin/env bash

grpcurl -plaintext -d '{"entities":[{"entity_id": 76565,"entity_type":"type1"}]}' \
    127.0.0.1:7678 evaluation.Evaluation/Evaluate
