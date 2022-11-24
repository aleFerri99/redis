# CoMoSyMe

[![N|Solid](https://d11wkw82a69pyn.cloudfront.net/concept-reply/en/siteassets/concept-reply.png)](https://nodesource.com/products/nsolid)

## _MODBUS Device Service_

MODBUS device is an EdgeX Device Service that can interact with MODBUS protocol, accordingly to the MODBUS configuration specified within a modbus.yaml configuration file: 

- It can handle MODBUS RTU, ASCII and TCP, accordingly to specified mode
- For purpose of this project only RTU use case has been taken into account
- ASCII end TCP are available but they have not been tested so far

## _modbus.yaml_ configuration file

The modbus configuration file is defined as follow:

```
sensors:
  <sensor_name>:
    name: t9501
    poll: true <ignored_for_container_case>
    mode: rtu <can_be_also_ascii_or_tcp>
    modbus-properties: # Are the same for ASCII/RTU case
      address: <slave_address>    # address: <tcp_address>
      baudrate: <slave_baud_rate> # port: <tcp_port>
      databits: <slave_data_bits> # //
      stopbits: <slave_stop_bits> # //
      parity: <slave_parity_bit>  # //
      timeout: <serial_timeout>   # timeout: <tcp_timeout>
    modbus:
      registers: <list_of_device_register_to_start_read>
      - 0x1389 # Firmware version register
      - 0x138D # Relative humidity register
      - 0x138E # Temperature register
      - 0x0FA4 # Config flag register
      - 0x0FA5 # Slave address register
      - 0x0FA6 # Baud rate register
      - 0x0FA9 # Parity register
      read-length: <list_of_number_of_coils_to_read_from_start_register>
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      read-write: <list_of_operations_to_perform_on_registers>
      - function-code: 0x04 <modbus_function_code>
        value: <modbus_value_to_write_on_register>
        - 0x00 <defined_only_for_write_function_codes>
      - function-code: 0x04
        value: 
        - 0x00
      - function-code: 0x04
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
  ph:
    name: pH
    poll: true
    mode: rtu
    modbus-properties:
      address: 0x03
      baudrate: 9600
      databits: 8
      stopbits: 1
      parity: "N"
      timeout: 1
    modbus:
      registers:
      - 0x00 # TEMPERATURE register
      - 0x01 # PH register
      - 0x02 # PH_CALIBRATION_AD register
      - 0x20 # TEMP_CMP_STATUS register
      - 0x30 # PH_CALIB_RAW_AD0 register
      - 0x31 # PH_CALIB_RAW_AD1 register
      - 0x32 # PH_CALIB_RAW_AD2 register
      - 0x0200 # SLAVE_ADDR register
      - 0x0201 # BAUD_RATE register
      - 0x0202 # MODBUS_PROTOCOL register
      - 0x0203 # PARITY_BITS register
      - 0x0204 # DATA_BITS register
      - 0x0205 # STOP_BITS register
      - 0x0206 # SERIAL_RESPONSE_DELAY register
      - 0x0207 # SERIAL_ACTIVE_OUTPUT_TIME register
      read-length:
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      read-write:
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
  conductivity:
    name: conductivity
    poll: true
    mode: rtu
    modbus-properties:
      address: 0x02
      baudrate: 9600
      databits: 8
      stopbits: 1
      parity: "N"
      timeout: 1
    modbus:
      registers:
      - 0x00 # TEMPERATURE register
      - 0x02 # EC register
      - 0x03 # SALINITY register
      - 0x04 # TDS register
      - 0x08 # EC_RAW_AD register
      - 0x20 # TEMP_CMP_STATUS register
      - 0x22 # EC_TEMP_COFF register
      - 0x23 # SALINITY_COFF register
      - 0x24 # TDS_COFF register
      - 0x25 # ELECTRODE_CONST register
      - 0x30 # EC_CALIB_1413 register
      - 0x31 # EC_CALIB_12880 register
      - 0x0200 # SLAVE_ADDR register
      - 0x0201 # BAUD_RATE register
      - 0x0202 # MODBUS_PROTOCOL register
      - 0x0203 # PARITY_BITS register
      - 0x0204 # DATA_BITS register
      - 0x0205 # STOP_BITS register
      - 0x0206 # SERIAL_RESPONSE_DELAY register
      - 0x0207 # SERIAL_ACTIVE_OUTPUT_TIME register
      read-length:
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      - 1
      read-write:
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
      - function-code: 0x03
        value: 
        - 0x00
  density:
    name: density
    poll: true
    mode: rtu
    modbus-properties:
      address: 0x01
      baudrate: 9600
      databits: 8
      stopbits: 1
      parity: "N"
      timeout: 1
    modbus:
      registers:
      - 0x00 # All data
      read-length:
      - 6
      read-write:
      - function-code: 0x03
        value: 
        - 0x00
```

- Device service read this configuration file at service startup
- MODBUS configuration file is mounted on the container from host
- It is possible to specify a new configuration file and use it by changing the _modbusfile_ field in EdgeXIreland/values.yaml file
- The MODBUS protocol runs on RS485 serial standard, through a serial resource, defined by the field _serialresource_ in EdgeXIreland/values.yaml file

## _Application Services_

For this project, 4 different edgex application services has been developed to handle data conversion and export:

- ph-filter-mqtt-telegraf
- conductivity-filter-mqtt-telegraf
- density-filter-mqtt-telegraf
- t9501-filter-mqtt-telegraf

These services, takes data coming from EdgeX message bus (redis -> Defined within configuration.toml files of both DS and As) and filters them with the following flow:

- Filter data by Device Service name
- Convert filtered data to JSON
- Create sensor data structure and parse data within it
- Send parsed data to MQTT server

## Helm chart

All these service are deployed through a k8s helm chart, that introduces some environmental variables, defined under the _env_ field of the EdgeXIreland/values.yaml file, to make these services more flexible:

#### EdgeX ENV vars

See EdgeX documentation for these env variables [EdgeX Ireland](https://docs.edgexfoundry.org/2.0/)

- loggingenableremote
- securitysecretstore
- edgexregistry
- servicehost

#### Device Service ENV vars

- verbose: Flag variable used to print/hide debug services message
- serialresource: Serial resource used to start RS485 communication
- readgap: Delay between consecutive sensors
- modbusfile: Modbus configuration yaml file

#### Application Services ENV vars

- <sensor_name>: **Must** be the same specified within the modbus.yaml configration file, since it is used to take the snippet of interest for that specific sensor from the edgex messagebus. There must be a different env var for each sensor
- verbose: Flag variable used to print/hide debug services message
- protocol: MQTT protocol (tcp:// or ws://)
- broker: MQTT broker address
- port: MQTT broker port
- <application_client_id>: Client ID used to connect to broker. It **must** be unique for each application service
- user: MQTT broker connection username (specified in Broker ACL)
- password: MQTT broker connection password (specified in Broker ACL)
- qos: MQTT Quality Of Service
- retained: MQTT retained message flag
- autoreconnect: MQTT autoreconnect flag
- cleansession: MQTT clean session flag
- store: MQTT storage path if message persistency is enabled
- connectretry: MQTT connection retry attempts
- connectionretryinterval: MQTT interval between consecutive reconnection attempts
- tlsconnection: Service flag used to specify if TLS **must** be included in connection
- clientauthority: TLS CA certificate
- clientcertificate: TLS client certificate
- clientkey: TLS client key
- ignore: Service flag used to skip MQTT message publish
- <sensor_topic>: Topic on which each application service publish its data. It **MUST** be unique to avoi data overwritting

For more details about MQTT protocol take a look here: [MQTT](https://www.hivemq.com/hivemq/mqtt-broker/?utm_source=adwords&utm_campaign=&utm_term=mqtt%20broker%20hivemq&utm_medium=ppc&hsa_tgt=kwd-1477471991744&hsa_cam=17472918619&hsa_src=g&hsa_net=adwords&hsa_kw=mqtt%20broker%20hivemq&hsa_ad=603517458425&hsa_grp=138870236058&hsa_ver=3&hsa_acc=3585854406&hsa_mt=e&gclid=CjwKCAjw9NeXBhAMEiwAbaY4lj7gP0Rd7fSEW-udgNTqgtp8ZkHR9BvHv0fBz0KrI5ZTEIQfmQYoDBoCQl8QAvD_BwE)

## Installation

To configure a new environment, perform the following steps:

- Install k3s:

```sh
curl -sfL https://get.k3s.io | sh -
```

- Install helm:

```sh
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

- Install k9s:

```sh
curl -sS https://webinstall.dev/k9s | bash
```

Once done, you can deploy the application using helm:

```sh
cd CoMoSyMe/EdgeXIreland
# Check that k3s service is up and running
sudo systemctl status k3s.service
# Make kubeconfig file readable
sudo chmod +r /etc/rancher/k3s/k3s.yaml
# Create project namespace
kubectl create namespace comosyme
# Deploy application
helm install comosyme ./ --kubeconfig /etc/rancher/k3s/k3s.yaml
```