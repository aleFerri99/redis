# Build and import MODBUS Device Service
echo "##########################################################"
echo "Building device-modbus"
echo "##########################################################"
cd /home/docker/CoMoSyMe/device-modbus
make docker version=0.2.0-amd64 arch=amd64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
sudo docker save --output /home/docker/CoMoSyMe/images/device-modbus.tar gufiregistry.cloud.reply.eu/comosyme/device-modbus:0.2.0-amd64

echo "##########################################################"
echo "Building device-gpiod"
echo "##########################################################"
cd /home/docker/CoMoSyMe/device-gpiod
make docker version=0.2.0-amd64 arch=amd64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
sudo docker save --output /home/docker/CoMoSyMe/images/device-gpio.tar gufiregistry.cloud.reply.eu/comosyme/device-gpiod:0.2.0-amd64

echo "##########################################################"
echo "Building device-mqtt"
echo "##########################################################"
cd /home/docker/CoMoSyMe/device-mqtt
make docker version=0.2.0-amd64 arch=amd64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
sudo docker save --output /home/docker/CoMoSyMe/images/device-mqtt.tar gufiregistry.cloud.reply.eu/comosyme/device-mqtt:0.2.0-amd64

# Build and import pH Service
echo "##########################################################"
echo "Building app-ph-mqtt-telegraf"
echo "##########################################################"
cd /home/docker/CoMoSyMe/ph-filter-mqtt-telegraf
make docker version=0.2.0-amd64 arch=amd64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
sudo docker save --output /home/docker/CoMoSyMe/images/app-ph.tar gufiregistry.cloud.reply.eu/comosyme/app-ph-mqtt-telegraf:0.2.0-amd64

# Build and import conductivity Service
echo "##########################################################"
echo "Building app-conductivity-mqtt-telegraf"
echo "##########################################################"
cd /home/docker/CoMoSyMe/conductivity-filter-mqtt-telegraf
make docker version=0.2.0-amd64 arch=amd64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
sudo docker save --output /home/docker/CoMoSyMe/images/app-conductivity.tar gufiregistry.cloud.reply.eu/comosyme/app-conductivity-mqtt-telegraf:0.2.0-amd64

# Build and import density Service
echo "##########################################################"
echo "Building app-density-mqtt-telegraf"
echo "##########################################################"
cd /home/docker/CoMoSyMe/density-filter-mqtt-telegraf
make docker version=0.2.0-amd64 arch=amd64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
sudo docker save --output /home/docker/CoMoSyMe/images/app-density.tar gufiregistry.cloud.reply.eu/comosyme/app-density-mqtt-telegraf:0.2.0-amd64

# Build and import concentration Service
echo "##########################################################"
echo "Building app-concentration-mqtt-telegraf"
echo "##########################################################"
cd /home/docker/CoMoSyMe/concentration-filter-mqtt-telegraf
make docker version=0.2.0-amd64 arch=amd64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
sudo docker save --output /home/docker/CoMoSyMe/images/app-concentration.tar gufiregistry.cloud.reply.eu/comosyme/app-concentration-mqtt-telegraf:0.2.0-amd64

# Build and import conductivity Service
echo "##########################################################"
echo "Building app-t9501-mqtt-telegraf"
echo "##########################################################"
cd /home/docker/CoMoSyMe/t9501-filter-mqtt-telegraf
make docker version=0.2.0-amd64 arch=amd64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
sudo docker save --output /home/docker/CoMoSyMe/images/app-t9501.tar gufiregistry.cloud.reply.eu/comosyme/app-t9501-mqtt-telegraf:0.2.0-amd64

# Build and import gpio Service
echo "##########################################################"
echo "Building app-gpio-mqtt-telegraf"
echo "##########################################################"
cd /home/docker/CoMoSyMe/gpio-filter-mqtt-telegraf
make docker version=0.2.0-amd64 arch=amd64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
sudo docker save --output /home/docker/CoMoSyMe/images/app-gpio.tar gufiregistry.cloud.reply.eu/comosyme/app-gpio-mqtt-telegraf:0.2.0-amd64

sudo chmod 777 /home/docker/CoMoSyMe/images/*.tar
sudo k3s ctr images import /home/docker/CoMoSyMe/images/device-gpio.tar
sudo k3s ctr images import /home/docker/CoMoSyMe/images/device-modbus.tar
sudo k3s ctr images import /home/docker/CoMoSyMe/images/device-mqtt.tar
sudo k3s ctr images import /home/docker/CoMoSyMe/images/app-gpio.tar
sudo k3s ctr images import /home/docker/CoMoSyMe/images/app-ph.tar
sudo k3s ctr images import /home/docker/CoMoSyMe/images/app-density.tar
sudo k3s ctr images import /home/docker/CoMoSyMe/images/app-t9501.tar
sudo k3s ctr images import /home/docker/CoMoSyMe/images/app-conductivity.tar
sudo k3s ctr images import /home/docker/CoMoSyMe/images/app-concentration.tar
