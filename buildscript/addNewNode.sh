#! /bin/bash

echo "addNewNode.sh localID(local IP) serverIP"

echo ""
echo "Step 0. Upgrade docker version"
curl -sS https://gist.githubusercontent.com/andy-zhangtao/69a947fc992c55828582148582d085c6/raw/d4610d226e679341842fcf3ec872d866253aceeb/upgrade-docker-ce.sh|bash

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

cd /usr/local/bin ; sftp yooadmin@182.254.217.30 <<EOT
get /tmp/log_watcher_slave
EOT

cat << EOT > /usr/local/bin/start_log_watcher_slave
#!/bin/bash
cd /usr/local/bin ; ./log_watcher_slave -vvv d -c $2 -i $1
EOT

echo "       /usr/local/bin/start_log_watcher_slave install finish"

chmod 755 /usr/local/bin/start_log_watcher_slave

cat << EOT > /etc/systemd/system/log_watcher_slave.service
[Unit]
Description=log_watcher_slave

[Service]
TimeoutStartSec=0

Restart=always
ExecStart=/usr/local/bin/start_log_watcher_slave
ExecStop=-/usr/bin/mv /tmp/log_watcher_slave /usr/local/bin/log_watcher_slave

[Install]
WantedBy=multi-user.target
EOT

systemctl daemon-reload
systemctl start log_watcher_slave.service

echo ""
echo "      /etc/systemd/system/log_watcher_slave.service install finish"

journalctl -fu log_watcher_slave

echo ""
echo "Modify openfile limit"
echo "modify /etc/security/limits.conf "
echo "* soft nofile 102400"
echo "* hard nofile 102400"

