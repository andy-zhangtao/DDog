#!/bin/bash
echo "sudo useradd yooadmin"
sudo useradd yooadmin
sudo su << EOF
adduser yooadmin sudo
mkdir -p /home/yooadmin
chown -R yooadmin /home/yooadmin 
su - yooadmin
mkdir -p /home/yooadmin/.ssh
echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDm5wFzJbopJLrSnAkZ6BsoY9BxlijA1eWOAcI45EJgUFaMaJgm5XGBWlxOMpcekkUu2eRyN6hNXAAScUPkrmXoemSq8gHCTZ0xa1bR4Jhx04RZY7JI/naKuRUP/sw1U82mWU8htydOBNbOtihzTe2+bn3aKa78dgUpVfQvk90FwAaC6LvPPbDrLKSws9OhfxAOE1uug4BbzX9ir1r0jfwzpq2JRq7GGtvSoAAR93UlkdHtpSdY46AX6nJyDO333TDvjWwY5NO439fSrEkGaVmf9NL82CDumBuZTIP/kAC/kR0z9/cq8Cp5I/6BeX4T619Ouvn8yjTQFawntpvUjjR9 jumpserver@279117151cdc" > /home/yooadmin/.ssh/authorized_keys
EOF
echo "create authorized_keys"
echo "chmod yooadmin passwd"
sudo passwd yooadmin << EOF
yooadmin
yooadmin
EOF

echo "docker pull vikings/gotty"
sudo docker pull vikings/gotty

echo "create log_watcher_slave.service"
sudo bash -c 'cat << EOF > /etc/systemd/system/log_watcher_slave.service
[Unit]
Description=log_watcher_slave

[Service]
TimeoutStartSec=0

Restart=always
ExecStart=/usr/local/bin/log_watcher_slave -vvv d -c 192.168.0.15:8000 -i 192.168.2.11
ExecStopPost=-/usr/bin/wget 192.168.0.15:8005/log_watcher_slave -O /usr/local/bin/log_watcher_slave
ExecStopPost=-/bin/chmod 755 /usr/local/bin/log_watcher_slave

[Install]
WantedBy=multi-user.target
EOF'

sudo systemctl daemon-reload
sudo systemctl start log_watcher_slave.service
echo ""
echo "      /etc/systemd/system/log_watcher_slave.service install finish"

sudo journalctl -fu log_watcher_slave