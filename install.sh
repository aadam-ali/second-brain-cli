#!/usr/bin/env bash

install_dir="${HOME}/.local/bin"

[[ -d "${install_dir}" ]] || mkdir -p "${install_dir}"

if [[ -f "${install_dir}/sb" ]]; then
  mv "${install_dir}/sb" "${install_dir}/sb.old"
fi

go get .
go build -o "${install_dir}/sb" -ldflags "-X github.com/aadam-ali/second-brain-cli/config.version=$(git rev-parse --short HEAD)" main.go
