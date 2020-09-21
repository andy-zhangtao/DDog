#!/bin/sh

/ddog & sleep 1 & nohup /coredns  $* >> /dns.log & tail -f /dns.log
