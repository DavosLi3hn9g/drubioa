#!/bin/bash
WORKSPACE="/home/pi/VGO"
cd ${WORKSPACE}
sudo ./iqiar $*

# sudo cp -f /home/pi/VGO/start.sh /etc/init.d/
# sudo chmod +x /etc/init.d/start.sh
# sudo chown root:root /etc/init.d/start.sh