#! /bin/bash
# export CHAINPATHS=("/media/disk1/chaindata/ethmainnet" "/media/disk3/chaindata1" "/media/disk3/chaindata2" "/media/disk4/chaindata3" "/media/disk4/chaindata4" "/media/disk1/chaindata/chaindata5" "/media/disk3/chaindata6" "/media/disk4/chaindata7")
export CHAINPATHS=("/media/disk1/chaindata/ethmainnet" "/media/disk3/chaindata1" "/media/disk3/chaindata2" "/media/disk4/chaindata3" "/media/disk4/chaindata4")
# export CHAINPATHS=("/media/disk1/chaindata/ethmainnet")
export RSYS_CONFIG=/etc/rsyslog.d/10-geth.conf
export LOGROTATE_CONFIG=/etc/logrotate.d/geth
export CRON_CONFIG=/etc/cron.d/geth
export GETH=/media/disk1/workstation/eth-nodes/go-ethereum/build/bin/geth
export PRYSM=/media/disk1/workstation/eth-nodes/prysm/prysm.sh
export SCRIPTPATH="$(
    cd -- "$(dirname "$0")" >/dev/null 2>&1
    pwd -P
)"

info() {
    echo -e "\033[1;32m---> ${*}\033[0m"
}

warn() {
    echo -e "\033[1;33m---> ${*}\033[0m"
}
error() {
    echo -e "\033[1;31m---> ${*}\033[0m"
}
