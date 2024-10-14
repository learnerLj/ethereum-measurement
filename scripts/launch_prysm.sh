#!/bin/bash
set -eux

cd "$(dirname "$(realpath "$0")")"
source "env.sh"
launch_prysm() {
    echo
    for i in "${!CHAINPATHS[@]}"; do
        local prysmsessname=prysm"$((i + 1))"
        info "launch $prysmsessname..."
        if tmux has-session -t "$prysmsessname" 2>/dev/null; then
            tmux send-keys -t "$prysmsessname" C-c C-m
        else
            tmux new-session -d -s "$prysmsessname"
        fi
        sleep 1
        tmux send-keys -t "$prysmsessname" "cd ${CHAINPATHS[i]}" C-m
        tmux send-keys -t "$prysmsessname" "$PRYSM beacon-chain --config-file=./prysm.yaml" C-m

    done
}
launch_prysm
