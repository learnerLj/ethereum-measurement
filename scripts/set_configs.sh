#!/bin/bash
set -eux

cd "$(dirname "$(realpath "$0")")"
source "env.sh"
set_configs() {
    cd "$(dirname "$0")"
    for i in "${!CHAINPATHS[@]}"; do
        echo
        local GETH_CONFIG="${CHAINPATHS[i]}/geth-config.toml"
        local PRYSM_CONFIG="${CHAINPATHS[i]}/prysm.yaml"

        info "generate jwt key in ${CHAINPATHS[i]}..."
        openssl rand -hex 32 | tr -d "\n" >"${CHAINPATHS[i]}/jwt.hex"

        info "set geth configuration in ${CHAINPATHS[i]}..."
        cp geth-config.toml "$GETH_CONFIG"

        awk -v i="$i" -v chainpath="${CHAINPATHS[i]}" -v jwt_path="${CHAINPATHS[i]}/jwt.hex" '
        $1 == "HTTPPort" { $3 = 8540 + i }
        $1 == "AuthPort" { $3 = 8600 + i }
        $1 == "ListenAddr" { $3 = "\":303" sprintf("%02d", 0 + i) "\"" }
        $1 == "DiscAddr" { $3 = "\":303" sprintf("%02d", 0 + i) "\"" }
        $1 == "DataDir" { $3 = "\"" chainpath "/el\"" }
        $1 == "JWTSecret" { $3 = "\"" jwt_path "\"" }
        { print }
    ' "$GETH_CONFIG" >tmp.toml && mv tmp.toml "$GETH_CONFIG"

        info "set prysm configuration in ${CHAINPATHS[i]}..."
        cp prysm.yaml "$PRYSM_CONFIG"
        awk -v i="$i" -v cl_path="${CHAINPATHS[i]}/cl" -v jwt_path="${CHAINPATHS[i]}/jwt.hex" '
        $1 == "datadir:" { $2 = "\"" cl_path "\"" }
        $1 == "jwt-secret:" { $2 = "\"" jwt_path "\"" }
        $1 == "execution-endpoint:" { $2 = "\"http://localhost:" 8600 + i "\"" }
        $1 == "p2p-tcp-port:" { $2 = 13000 + i }
        $1 == "p2p-udp-port:" { $2 = 12000 + i }
        $1 == "grpc-gateway-port:" { $2 = 3500 + i }
        $1 == "monitoring-port:" { $2 = 8080 + i }
        $1 == "rpc-port:" { $2 = 4000 + i }
        { print }
    ' "$PRYSM_CONFIG" >tmp.yaml && mv tmp.yaml "$PRYSM_CONFIG"
    done
}
set_configs
