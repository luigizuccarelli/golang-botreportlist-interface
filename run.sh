#!/bin/bash

LOG_LEVEL=debug
SERVER_PORT=9000
VERSION=1.0.1
NAME=servisbot-lytics-interface
URL=https://api.lytics.io/collect/json/user_db?access_token=a2Djtb6hna3biMbgaOXgAQxx

export NAME LOG_LEVEL SERVER_PORT VERSION URL

./build/microservice
