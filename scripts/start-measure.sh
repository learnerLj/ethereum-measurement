#!/bin/bash
set -eux

cd "$(dirname "$(realpath "$0")")"

source "env.sh"
main() {
  if [[ $EUID -ne 0 ]]; then
    error "This script must be run as root"
    exit 1
  fi
  initiate_logger

  sudo -u "$SUDO_USER" "$SCRIPTPATH"/set_configs.sh
  sudo -u "$SUDO_USER" "$SCRIPTPATH"/launch_geth.sh
  sudo -u "$SUDO_USER" "$SCRIPTPATH"/launch_prysm.sh
  info "SUCCESS!!!!"
}

initiate_logger() {
  apt-get install -y logrotate rsyslog cron
  rm "$RSYS_CONFIG" "$LOGROTATE_CONFIG" "$CRON_CONFIG"

  info "make directors..."

  for chainpath in "${CHAINPATHS[@]}"; do
    mkdir -p "$chainpath"/{el,cl} && chown -R "$SUDO_USER":"$SUDO_USER" "$chainpath"
    sudo -u "$SUDO_USER" mkdir -p "$chainpath"/{el,cl}
    cd "$chainpath"
    touch geth.log && chown "$SUDO_USER":syslog geth.log && chmod 664 geth.log
  done

  info "config rsyslog..."
  if [ ! -f $RSYS_CONFIG ]; then
    printf "template(name=\"PlainMsg\" type=\"string\" string=\"%%msg%%\\\\n\")\n" >"$RSYS_CONFIG"
    for i in "${!CHAINPATHS[@]}"; do
      {
        printf "\nif \$programname == 'geth%d' then {\n" "$((i + 1))"
        printf "action(type=\"omfile\" file=\"%s/geth.log\" template=\"PlainMsg\")\n" "${CHAINPATHS[i]}"
        printf "stop\n}\n"
      } >>"$RSYS_CONFIG"
    done
    chmod 644 $RSYS_CONFIG
  else
    warn "rsyslog config of geth exists"
  fi

  info "config logrotate..."
  if [ ! -f $LOGROTATE_CONFIG ]; then
    for i in "${!CHAINPATHS[@]}"; do
      {
        printf "%s/geth.log {\n" "${CHAINPATHS[i]}"
        printf "    size 100M\n"
        printf "    rotate 20\n"
        printf "    compress\n"
        printf "    missingok\n"
        printf "    notifempty\n"
        printf "    delaycompress\n"
        printf "    copytruncate\n"
        printf "    su %s syslog\n" "$SUDO_USER"
        printf "}\n"
      } >>"$LOGROTATE_CONFIG"
    done
    chmod 644 $LOGROTATE_CONFIG
  else
    warn "logrotate config of geth exists"
  fi

  info "config crontab..."
  if [ ! -f $CRON_CONFIG ]; then
    printf "*/2 * * * * root /usr/sbin/logrotate %s" "$LOGROTATE_CONFIG" >>$CRON_CONFIG
    chmod 644 $CRON_CONFIG
  else
    warn "cron config of geth exists"
  fi
  systemctl restart rsyslog logrotate cron
}

main
