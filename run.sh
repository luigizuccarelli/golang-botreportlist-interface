#!/bin/bash

LOG_LEVEL=debug
SERVER_PORT=9000
VERSION=1.0.1
NAME=servisbot-middleware-interface
URL=https://servicegateway.agora-inc.com/middleware/

export NAME LOG_LEVEL SERVER_PORT VERSION URL

./build/microservice
