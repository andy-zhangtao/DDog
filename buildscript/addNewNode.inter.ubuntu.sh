#! /bin/bash

echo "addNewNode.sh localID(local IP) serverIP"

echo ""
echo "Step 0. Upgrade docker version"
curl -sS https://gist.githubusercontent.com/andy-zhangtao/ea24ac8afdab0dd139afbfa01d39d81a/raw/59184d356b93732d378592fa660e586f8ab7940a/upgrade-docker-ce-ubuntu.sh|bash

echo "systemctl disable dockerd"
systemctl disable dockerd

echo "Step 0.5 Pull vikings/gotty"
docker pull vikings/gotty

         
echo ""
echo "Step 1. pull vikings/doctor-worker and copy /worker to /tmp/worker"
docker run -it --rm -v /tmp:/tmp --entrypoint="sh" vikings/doctor-worker
/tmp/worker -s root

echo ""
echo "Step 2. exec /tmp/worker -s yooadmin"
su - yooadmin <<!
/tmp/worker -s yooadmin
exit
!

echo ""
echo "Step 3. install log_watcher_slave service"


cat << EOT > /etc/systemd/system/log_watcher_slave.service
[Unit]
Description=log_watcher_slave

[Service]
TimeoutStartSec=0

Restart=always
ExecStart=/usr/local/bin/log_watcher_slave -vvv d -c 192.168.0.15:8000 -i xxxxxxxxxx
ExecStopPost=-/usr/bin/wget 192.168.0.15:8005/log_watcher_slave -O /usr/local/bin/log_watcher_slave
ExecStopPost=-/bin/chmod 755 /usr/local/bin/log_watcher_slave

[Install]
WantedBy=multi-user.target
EOT

systemctl daemon-reload
systemctl start log_watcher_slave.service

echo ""
echo "      /etc/systemd/system/log_watcher_slave.service install finish"

journalctl -fu log_watcher_slave