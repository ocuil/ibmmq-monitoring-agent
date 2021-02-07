# ibmmq-monitoring-agent to elastic

- This code is base on [IBM® MQ golang metric packages](https://github.com/ibm-messaging/mq-golang) official repository, __so all code It's not my own__ :exclamation:
- The motive to fork it, it is adapted to make an easy integration with [elasticsearch](https://www.elastic.co/es/elasticsearch/) and or [logstash](https://www.elastic.co/es/logstash) as an agent
- Fixed the problems __cmqc.h No such file or directory__ :fire: and similar, caused by the absence of c headers in the import sections, go.mod and vendors
- The laboratory we use to adapt and finally build the agent is Ubuntu 18.04.5 LTS, but changing to the appropriate command line yum and a right way to install golang, you'll should don't get a problem with that.
- It's not ready yet to work with elasticsearch. __work in progress__

## Important
![Golang](https://img.shields.io/badge/Go-1.15.7-blue)
![Ubuntu](https://img.shields.io/badge/ubuntu-18.04.5%20LTS-red)
![Version](https://raw.githubusercontent.com/ocuil/assets/main/img/version.svg)
![Release](https://img.shields.io/badge/release-alpha-brightgreen)
![Build](https://img.shields.io/badge/build-passing-brightgreen)
![Implementations](https://raw.githubusercontent.com/ocuil/assets/main/img/non-for-production.svg)
![Size](https://img.shields.io/github/languages/code-size/ocuil/ibmmq-monitoring-agent)
![RepoSize](https://img.shields.io/github/repo-size/ocuil/ibmmq-monitoring-agent)
![Lines](https://img.shields.io/tokei/lines/github/ocuil/ibmmq-monitoring-agent)

---
# asciicast video about laboratory
---

[![asciicast](https://asciinema.org/a/Mjt2Uco4nTmYfHqYHuKS4ICiD.svg)](https://asciinema.org/a/Mjt2Uco4nTmYfHqYHuKS4ICiD)

---
---
# Install build tools (ubuntu) and download repo
```bash
sudo apt update -y && sudo apt upgrade -y && sudo apt install build-essential -y
#upload the IBM-MQC
git clone https://github.com/ocuil/ibmmq-monitoring-agent.git
wget https://golang.org/dl/go1.15.7.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.7.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

---
---
# Install MQ Client
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

# Build the binary 

```bash
#inside the project folder
env GOOS=linux GOARCH=amd64 go build -mod=vendor ./...
```
---
---
# Install IBM MQ to create lab env and play
__I will work to improve the documentation__

_first_ you have to delete `/opt/mqm` if you (like me) are recycling the machine, thats because the IBM MQ install process would like to use that folder, possible its can be change but I have no experience with it software so I'll try to keep all closer to the default. 

- Before you have to download the software from the IBM offical portal, you can register or use your IBM ID to get a Trial version.

__As I write before, the lab is base on Ubuntu so I use the Ubuntu compatible packages to do this__

- The order to install packages trought `dpkg` in the official site of IBM are not really ready, because put the `ibmmq-serve` installed step before `ibmmq-gskit`, and its his dependency
- You can use the `apt` if you want, it is explained by IBM in his documentation too

```bash
dpkg -i ibmmq-runtime_9.2.0.0_amd64.deb
dpkg -i ibmmq-jre_9.2.0.0_amd64.deb
dpkg -i ibmmq-java_9.2.0.0_amd64.deb
dpkg -i ibmmq-gskit_9.2.0.0_amd64.deb
dpkg -i ibmmq-server_9.2.0.0_amd64.deb
dpkg -i ibmmq-server_9.2.0.0_amd64.deb
#Here you will get a warning about the recommendations... but it's a lab no a prodction install =)
dpkg -i ibmmq-web_9.2.0.0_amd64.deb
dpkg -i ibmmq-ftbase_9.2.0.0_amd64.deb
dpkg -i ibmmq-ftagent_9.2.0.0_amd64.deb
dpkg -i ibmmq-ftservice_9.2.0.0_amd64.deb
dpkg -i ibmmq-ftlogger_9.2.0.0_amd64.deb
dpkg -i ibmmq-fttools_9.2.0.0_amd64.deb
dpkg -i ibmmq-amqp_9.2.0.0_amd64.deb
dpkg -i ibmmq-ams_9.2.0.0_amd64.deb
dpkg -i ibmmq-xrservice_9.2.0.0_amd64.deb
dpkg -i ibmmq-explorer_9.2.0.0_amd64.deb
dpkg -i ibmmq-client_9.2.0.0_amd64.deb
dpkg -i ibmmq-man_9.2.0.0_amd64.deb
dpkg -i ibmmq-msg-es_9.2.0.0_amd64.deb
dpkg -i ibmmq-samples_9.2.0.0_amd64.deb
dpkg -i ibmmq-sdk_9.2.0.0_amd64.deb
dpkg -i ibmmq-sfbridge_9.2.0.0_amd64.deb
dpkg -i ibmmq-bcbridge_9.2.0.0_amd64.deb
```

## Time to setup the IBM MQ to use the agent:

- Create the manager => ```crtmqm -q gravity```
- Start the Manager => ```strmqm gravity```
- Get into the console to create the queue => ```runmqsc gravity```
- Create the queue => ```define ql(gravity.cola01)```
- Check the queue => ```dspmq```
- Add a message => ```/opt/mqm/samp/bin/amqsput GRAVITY.COLA01``` (2 empty enters to finish)

## Now is the moment to config the shell script 'mq_json.sh'
- The manager you will get metrics ```queues="GRAVITY.*"```
- The path to the agent ```exec /home/mqm/ibmmq-monitoring-agent $ARGS```

## Create the service on IBM MQ that will execute the shell script that execute the agent

```
DEFINE SERVICE(MQJSON)         +
       CONTROL(QMGR)               +
       SERVTYPE(SERVER)            +
       STARTCMD('/home/mqm/mq_json.sh') +
       STARTARG(+QMNAME+)          +
       STOPCMD('/usr/bin/kill -9' )  +
       STOPARG(+MQ_SERVER_PID+)    +
       STDOUT('/var/mqm/errors/mq_json.out')  +
       STDERR('/var/mqm/errors/mq_json.err')  +
       DESCR('MQ exporter for JSON format')
```

Check if the service are correctly setup ```DISPLAY SVSTATUS(MQJSON)```

Start the service ```START SERVICE(MQJSON)```

If you want to start again:

```
STOP SERVICE(MQJSON)
DELETE SERVICE(MQJSON)
```

---
---
---
# Always made with passion on golang !!!
<p align="center">
<img src="https://raw.githubusercontent.com/ocuil/assets/main/img/heart-hug.svg">
</p>
