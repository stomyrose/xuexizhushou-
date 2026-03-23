#!/bin/bash

echo "=== 查看服务日志 ==="

SERVICE=${1:-backend}

docker-compose logs -f $SERVICE
