#!/usr/bin/env bash
aws lambda invoke --function-name gym_occupancy response.json --cli-binary-format raw-in-base64-out

cat response.json
