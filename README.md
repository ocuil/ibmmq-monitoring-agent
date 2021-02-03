# ibmmq-monitoring-agent to elastic

- This code is base on [IBM MQ](https://github.com/ibm-messaging/mq-metric-samples) official repository
- The motive to fork it is adatpt to make an easy integration with [elasticsearch](https://www.elastic.co/es/elasticsearch/) and or [logstash](https://www.elastic.co/es/logstash)
- Fixed the problems __cmqc.h No such file or directory__ and similar, caused by the absence of c headers

## Important
![Version](https://raw.githubusercontent.com/ocuil/assets/main/img/version.svg)
![GitHub Release](https://raw.githubusercontent.com/ocuil/assets/main/img/release.svg)
![Implementations](https://raw.githubusercontent.com/ocuil/assets/main/img/non-for-production.svg)

---
---
# Patch the vendors source files to add the c headers files that need go to compile the agent
### files that added to the vendor folder:
```
vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/includes
 · cmqc.h
 · cmqcfc.h
 · cmqstrc.h
 · cmqxc.h
```
### files that are patched:
```
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqi.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiBO.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiCBC.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiCBD.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiCTLO.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiDLH.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMHO.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMPO.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMQCD.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMQCNO.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMQGMO.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMQMD.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMQOD.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMQPMO.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMQSCO.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiMQSD.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiPCF.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiSRO.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiSTS.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqiattrs.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqicb.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqicb_c.go
modified:   vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/mqistr.go

```
---
---
# Install MQ Client
## Notes (maybe it's not necessary if you update the vendor files to add the c headers that need to be compiled for agent)
The lab env are in Ubuntu 18.04.5 LTS and to compile go initially we need to install "IBM-MQC-Redist" so after download it from the oficial website of IBM in this version 9.1.5.o ("9.1.5.0-IBM-MQC-Redist-LinuxX64.tar.gz") we perform the next lines:
```bash
mkdir /opt/mqm

#copy the content of tar.gz in dis folder

genmqpkg_incnls=1 \
genmqpkg_incsdk=1 \
genmqpkg_inctls=1

cd /opt/mqm
sudo bin/genmqpkg.sh -b /opt/mqm
```

output:

```
Generate MQ Runtime Package
---------------------------
This program will help determine a minimal set of runtime files that are
required for a queue manager installation or to be distributed with a
client application. The program will ask a series of questions and then
prompt for a filesystem location for the runtime files.

Note that IBM can only provide support assistance for an unmodified set
of runtime files.


The MQ runtime package will be created in

/opt/mqm


Generation complete !
MQ runtime package created in '/opt/mqm'
```

Config vars:
```bash
LD_LIBRARY_PATH="/opt/mqm/lib64:/usr/lib64" \
MQ_CONNECT_TYPE=CLIEN
```