# Stack of telegraf, alloy, loki, grafana
## How to start the stack
- boot up the compose file with docker compose
- send testtrap with testtrapsender ./testtrapsender/trapsenderv1_XX_XX use the build for your os. Flags can be used.
- navigate to 127.0.0.1:3000 to check the logentry in grafanaUI

## Check Configs/Files of Stack Apps
- check alloy config ./alloy/alloy.config
- check telegraf config ./telegraf/telegraf.conf
- check mibs dir ./mibs

## PreConfigured
### Telegraf
- able to receive and translate snmpv1traps
- able to output translated trap to loki api (here alloy source api component)

### Alloy
- able to receive loki reqs
- able to relabel and process the logentry
- able to write it to actual loki

### loki
- basic local install (no persistence)

### grafana
- basic local install with preconfigured loki datasource
