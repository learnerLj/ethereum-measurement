#!/bin/bash
set -eux

cd "$(dirname "$(realpath "$0")")"
source "env.sh"
launch_geth() {
    echo
    for i in "${!CHAINPATHS[@]}"; do
        local gethsessname=geth"$((i + 1))"
        info "launch $gethsessname..."
        if tmux has-session -t "$gethsessname" 2>/dev/null; then
            tmux send-keys -t "$gethsessname" C-c C-m
        else
            tmux new-session -d -s "$gethsessname"
        fi
        sleep 1
        tmux send-keys -t "$gethsessname" "cd ${CHAINPATHS[i]}" C-m
        tmux send-keys -t "$gethsessname" "$GETH --config geth-config.toml --verbosity 3 2>&1 | logger -t $gethsessname" C-m
    done
}

launch_geth
