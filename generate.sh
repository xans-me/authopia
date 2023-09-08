#!/bin/bash
protoc -I ./core/proto \
  --go_out ./core/proto --go_opt paths=source_relative \
  --go-grpc_out ./core/proto --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./core/proto --grpc-gateway_opt paths=source_relative \
  --openapiv2_out ./core/proto --openapiv2_opt use_go_templates=true \
  ./core/proto/*.proto
