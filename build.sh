#!/usr/bin/env bash
set -ex

platforms=("darwin/arm64" "windows/amd64" "darwin/amd64" "linux/amd64" "linux/arm64")

mkdir -p build
for platform in "${platforms[@]}"; do
    rm -f build/*
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    go build -o build/
    upx build/*
    output=$(ls build)
    tgz_name=${output%.*}'-'$GOOS'-'$GOARCH'.tgz'
    tar -C build -czf $tgz_name $output
done

