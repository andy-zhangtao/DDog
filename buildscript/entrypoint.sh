#!/bin/sh
if [[ -z "${RUN_TIME_DNS}" ]]; then
    echo "RUN_TIME_DNS not set. Use default Dns nameserver"
else
    echo nameserver $RUN_TIME_DNS >/etc/resolv.conf
fi
cat /etc/resolv.conf
