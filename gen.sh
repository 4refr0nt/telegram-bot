#!/bin/bash
swagger generate server -f api/api.yml -A "telegram-bot" --server-package=internal/handler 
