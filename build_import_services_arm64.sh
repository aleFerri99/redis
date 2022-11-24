# Build and import MODBUS Device Service
echo "##########################################################"
echo "Building device-modbus"
echo "##########################################################"
cd device-modbus
make docker version=0.2.0-arm64 arch=arm64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
docker push gufiregistry.cloud.reply.eu/comosyme/device-modbus:0.2.0-arm64

echo "##########################################################"
echo "Building device-gpiod"
echo "##########################################################"
cd device-gpiod
make docker version=0.2.0-arm64 arch=arm64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
docker push gufiregistry.cloud.reply.eu/comosyme/device-gpiod:0.2.0-arm64

echo "##########################################################"
echo "Building device-mqtt"
echo "##########################################################"
cd device-mqtt
make docker version=0.2.0-arm64 arch=arm64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
docker push gufiregistry.cloud.reply.eu/comosyme/device-mqtt:0.2.0-arm64

# Build and import pH Service
echo "##########################################################"
echo "Building app-ph-mqtt-telegraf"
echo "##########################################################"
cd ph-filter-mqtt-telegraf
make docker version=0.2.0-arm64 arch=arm64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
docker push gufiregistry.cloud.reply.eu/comosyme/app-ph-mqtt-telegraf:0.2.0-arm64

# Build and import conductivity Service
echo "##########################################################"
echo "Building app-conductivity-mqtt-telegraf"
echo "##########################################################"
cd conductivity-filter-mqtt-telegraf
make docker version=0.2.0-arm64 arch=arm64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
docker push gufiregistry.cloud.reply.eu/comosyme/app-conductivity-mqtt-telegraf:0.2.0-arm64

# Build and import concentration Service
echo "##########################################################"
echo "Building app-concentration-mqtt-telegraf"
echo "##########################################################"
cd concentration-filter-mqtt-telegraf
make docker version=0.2.0-arm64 arch=arm64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
docker push gufiregistry.cloud.reply.eu/comosyme/app-concentration-mqtt-telegraf:0.2.0-arm64

# Build and import density Service
echo "##########################################################"
echo "Building app-density-mqtt-telegraf"
echo "##########################################################"
cd density-filter-mqtt-telegraf
make docker version=0.2.0-arm64 arch=arm64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
docker push gufiregistry.cloud.reply.eu/comosyme/app-density-mqtt-telegraf:0.2.0-arm64

# Build and import t9501 Service
echo "##########################################################"
echo "Building app-t9501-mqtt-telegraf"
echo "##########################################################"
cd t9501-filter-mqtt-telegraf
make docker version=0.2.0-arm64 arch=arm64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
docker push gufiregistry.cloud.reply.eu/comosyme/app-t9501-mqtt-telegraf:0.2.0-arm64

# Build and import gpio Service
echo "##########################################################"
echo "Building app-gpio-mqtt-telegraf"
echo "##########################################################"
cd gpio-filter-mqtt-telegraf
make docker version=0.2.0-arm64 arch=arm64
echo "##########################################################"
echo "Saving docker image"
echo "##########################################################"
docker push gufiregistry.cloud.reply.eu/comosyme/app-gpio-mqtt-telegraf:0.2.0-arm64
