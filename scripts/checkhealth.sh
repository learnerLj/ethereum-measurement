#!/bin/bash
set -eu
cd "$(dirname "$(realpath "$0")")"
source "env.sh"

check_health() {
    local cmd1="admin.peers.length"
    local cmd2="eth.blockNumber"
    local cmd3="txpool.status"
    local cmd4=("admin.peers.reduce((acc, {name}) => (acc[name.split('/')[0]] = (acc[name.split('/')[0]] || 0) + 1, acc), {})")

    local total_peers=0
    local geth=$GETH
    unset GETH

    # 创建一个临时文件用于存储客户端类型和计数
    local client_types_file=$(mktemp)

    for i in "${!CHAINPATHS[@]}"; do
        local ipc="${CHAINPATHS[i]}/el/geth.ipc"
        info "geth$((i + 1)) status:"

        local block_number=$($geth attach --exec $cmd2 "$ipc")
        printf "%-20s: %s\n" "BlockNumber" "$block_number"

        local mem_pool_json=$($geth attach --exec $cmd3 "$ipc" | sed 's/\([a-zA-Z]\+\):/"\1":/g')
        local pending=$(echo "$mem_pool_json" | jq '.pending')
        local queued=$(echo "$mem_pool_json" | jq '.queued')
        printf "%-20s: %s\n" "Txpool.pending" "$pending"
        printf "%-20s: %s\n" "Txpool.queued" "$queued"

        local peers=$($geth attach --exec $cmd1 "$ipc")
        printf "%-20s: %s\n" "Peers" "$peers"
        total_peers=$((total_peers + peers))

        #    {: 1,Geth: 294} -> {: 1,"Geth": 294} -> {"Others" : 1,"Geth": 294}
        local peer_info=$($geth attach --exec "${cmd4[@]}" "$ipc" | sed 's/\([a-zA-Z]\+\):/"\1":/g' | sed 's/^\(\s*\):/\1"Others":/')
        echo "$peer_info" | jq -r 'to_entries[] | "\(.key) \(.value)"' | while read -r client count; do
            # 使用 awk 更新客户端类型的计数
            awk -v client="$client" -v count="$count" '
                BEGIN { found = 0; }
                $1 == client { $2 += count; found = 1; print; next; }
                { print; }
                END { if (!found) print client, count; }
            ' "$client_types_file" >"${client_types_file}.tmp" && mv "${client_types_file}.tmp" "$client_types_file"
        done

    done
    info "total peers: $total_peers"

    cat "$client_types_file" | while read -r client count; do
        percentage=$(echo "scale=2; 100 * $count / $total_peers" | bc)
        printf "%-20s: %-5s, %s%%\n" "$client" "$count" "$percentage"
    done

    # 删除临时文件
    rm "$client_types_file"
}

split_windows() {
    # tmux kill-session -t monitor
    # tmux kill-session -t admin
    # sleep 1

    # 创建新会话
    tmux new-session -d -s monitor
    tmux new-session -d -s admin
    for i in "${!CHAINPATHS[@]}"; do
        local log="${CHAINPATHS[i]}/geth.log"
        local ipc="${CHAINPATHS[i]}/el/geth.ipc"
        tmux send-keys -t monitor "cd ${CHAINPATHS[i]}" C-m
        tmux send-keys -t monitor "tail -f $log" C-m
        tmux split-window -v -t monitor

        tmux send-keys -t admin "cd ${CHAINPATHS[i]}" C-m
        tmux send-keys -t admin "$GETH attach $ipc" C-m
        tmux split-window -v -t admin

        tmux select-layout -t monitor even-vertical
        tmux select-layout -t admin even-vertical

    done
    tmux kill-pane -t monitor
    tmux kill-pane -t admin

    tmux select-layout -t monitor even-vertical
    tmux select-layout -t admin even-vertical
}

check_config() {
    for i in "${!CHAINPATHS[@]}"; do
        info "geth$((i + 1)):"
        local conf="${CHAINPATHS[i]}/geth-config.toml"
        # Check HTTPPort
        local http_port=$((8540 + i))
        grep "HTTPPort.*$http_port" "$conf" || echo "HTTPPort $http_port not found in $conf"

        # Check AuthPort
        local auth_port=$((8600 + i))
        grep "AuthPort.*$auth_port" "$conf" || echo "AuthPort $auth_port not found in $conf"

        # Check ListenAddr
        local listen_addr=":303$(printf "%02d" $((0 + i)))"
        grep "ListenAddr.*$listen_addr" "$conf" || echo "ListenAddr $listen_addr not found in $conf"

        # Check DiscAddr
        local disc_addr=":303$(printf "%02d" $((0 + i)))"
        grep "DiscAddr.*$disc_addr" "$conf" || echo "DiscAddr $disc_addr not found in $conf"

        info "prysm$((i + 1)):"
        local prysmconf="${CHAINPATHS[i]}/prysm.yaml"

        # Check execution-endpoint
        local execution_endpoint="http://localhost:$((8600 + i))"
        grep "execution-endpoint: \"$execution_endpoint\"" "$prysmconf" || echo "execution-endpoint $execution_endpoint not found in $prysmconf"

        # Check p2p-tcp-port
        local p2p_tcp_port=$((13000 + i))
        grep "p2p-tcp-port: $p2p_tcp_port" "$prysmconf" || echo "p2p-tcp-port $p2p_tcp_port not found in $prysmconf"

        # Check p2p-udp-port
        local p2p_udp_port=$((12000 + i))
        grep "p2p-udp-port: $p2p_udp_port" "$prysmconf" || echo "p2p-udp-port $p2p_udp_port not found in $prysmconf"

        # Check grpc-gateway-port
        local grpc_gateway_port=$((3500 + i))
        grep "grpc-gateway-port: $grpc_gateway_port" "$prysmconf" || echo "grpc-gateway-port $grpc_gateway_port not found in $prysmconf"

        # Check monitoring-port
        local monitoring_port=$((8080 + i))
        grep "monitoring-port: $monitoring_port" "$prysmconf" || echo "monitoring-port $monitoring_port not found in $prysmconf"

        # Check rpc-port
        local rpc_port=$((4000 + i))
        grep "rpc-port: $rpc_port" "$prysmconf" || echo "rpc-port $rpc_port not found in $prysmconf"

        echo
    done
}

main() {
    if [ $# -eq 0 ]; then
        warn "No arguments provided."
        check_health
        exit 0

    fi

    matched=false
    for arg in "$@"; do
        case $arg in
        windows)
            split_windows
            matched=true
            ;;
        health)
            check_health
            matched=true
            ;;
        config)
            check_config
            matched=true
            ;;
        esac
    done

    if [ "$matched" = false ]; then
        warn "no match command"
    fi

}

main "$@"
