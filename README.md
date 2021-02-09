# ibmmq-monitoring-agent to elastic


## disclaimer
- This code is base on [IBM® MQ golang metric packages](https://github.com/ibm-messaging/mq-golang) official repository, __so all code It's not my own__ :exclamation:

---

## why make it?

I am a passionate guy for elasticsearch and golang, I always try to learn, sometimes by myself (foros, chats, guides...) and sometimes with some [`#udemy`](#udemy) courses on my little free time.

So in a proyect where I am working, the team had the need to monitoring [`#IBM`](#IBM) MQ, [`#elastic`](#elastic) have in their ecosystem a module of [`#metricbeat`](metricbeat) &rightarrow; <a href="https://www.elastic.co/guide/en/beats/metricbeat/current/metricbeat-module-ibmmq.html" target="_blank">metricbeat-module-ibmmq</a>. But that module only works if the `IBMMQ` are installed `containerized`

<p align="center">
<img src="https://raw.githubusercontent.com/ocuil/assets/main/img/doc_elastic.png">
</p>

After read more deeply the documentation and check the git repo of `ibm-messaging` I was select the base of the actuall development agent.

`IBM` have an example that export to [`#JSON`](#JSON) over the standar sysout the metrics of the `MQ` &rightarrow; <a href="https://github.com/ibm-messaging/mq-metric-samples/tree/master/cmd/mq_json" target="_blank">MQ Exporter for JSON-based monitoring</a>

And as I writed before, this is an entirely `IBM's` code, just I made a few changes, maybe base on my ignorance. This changes add a C headers to a folder under the go files and change the include of them in the vendor version.

So that's it, we need metrics from `IBM MQ` services, and we have no integrations to do that with `elasticsearch`, then is why I start work on it

---

## current situation

![Golang](https://img.shields.io/badge/Go-1.15.7-blue)
![Ubuntu](https://img.shields.io/badge/ubuntu-18.04.5%20LTS-red)

![Version](https://raw.githubusercontent.com/ocuil/assets/main/img/version.svg)
![Release](https://img.shields.io/badge/release-alpha-brightgreen)
![Build](https://img.shields.io/badge/build-passing-brightgreen)
![Build](https://img.shields.io/badge/elasticsearch-ready-brightgreen)

![Size](https://img.shields.io/github/languages/code-size/ocuil/ibmmq-monitoring-agent)
![RepoSize](https://img.shields.io/github/repo-size/ocuil/ibmmq-monitoring-agent)
![Lines](https://img.shields.io/tokei/lines/github/ocuil/ibmmq-monitoring-agent)

A few hours before write that I finish to insert metrics from `IBM MQ` laboratory into `elasticsearch` sucesfully:
<p align="center">
<img src="https://raw.githubusercontent.com/ocuil/assets/main/img/discovery.png">
</p>

It was awesome to me, because it took me long long time (a lot of hours, weekends, nights ...), I love Golang but I don't use it on my day to day work, and I never work with `IBM MQ` before, thats was a handicap, looking for information, `how to` ...

---

## laboratory

After work with everything I noticed that we don't need to use the `IBM MQ` Service like they write on theirs github to run the agent, so you can choose the deep of the laboratory, I mean, just the need you to compile or if you want to play with

Just to compile the agent you need to install the `IBM MQ` Client and always patch the vendors folder (`patchVendors.sh`)
But if you want to play you can install just the `IBM MQ` Server and ... [`#have_fun`](#have_fun)

## __The `IBM MQ` Client and Server have to be downloaded previously from the IBM's website__ :exclamation: :exclamation:

---
## deploy a minimum laboratory &rightarrow; ![Golang](https://img.shields.io/badge/Go-1.15.7-blue) ![Ubuntu](https://img.shields.io/badge/ubuntu-18.04.5%20LTS-red)

In this case just we'll setup the `IBM MQ` Client, clone the repository, patch and compile it:

You can take a look to the next video:


[![asciicast](https://asciinema.org/a/Mjt2Uco4nTmYfHqYHuKS4ICiD.svg)](https://asciinema.org/a/Mjt2Uco4nTmYfHqYHuKS4ICiD)


### Step by step

1. Update and install build tools
2. Install the `IBM MQ` Client
3. Patch the `vendors`
4. Compile the agent


#### :one: Install build tools (ubuntu) and download repo
```bash
sudo apt update -y && sudo apt upgrade -y && sudo apt install build-essential -y
#upload the IBM-MQ previously download from IBM website
git clone https://github.com/ocuil/ibmmq-monitoring-agent.git
wget https://golang.org/dl/go1.15.7.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.7.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

#### :two: Install MQ Client
The lab env are in Ubuntu 18.04.5 LTS and to compile go initially we need to install "IBM-MQC-Redist" so after download it from the oficial website of IBM in this version 9.1.5.0 ("9.1.5.0-IBM-MQC-Redist-LinuxX64.tar.gz") we perform the next lines:
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

#### :three: Patch the vendors source files to add the c headers files that need go to compile the agent

The script `patchVendors.sh` perform that in a easy and fast way, just run `./patchVendors.sh`

```bash
#!/bin/bash
[ ! -d "./vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/includes" ] && cp -aR CHeaders ./vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/includes

find ./vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/ -type f -name "mqi*" -exec \
    sed -i -e 's%<cmqc.h>%"includes/cmqc.h"%g' {} +

find ./vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/ -type f -name "mqi*" -exec \
    sed -i -e 's%<cmqxc.h>%"includes/cmqxc.h"%g' {} +
    
find ./vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/ -type f -name "mqi*" -exec \
    sed -i -e 's%<cmqcfc.h>%"includes/cmqcfc.h"%g' {} +

find ./vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/ -type f -name "mqi*" -exec \
    sed -i -e 's%<cmqstrc.h>%"includes/cmqstrc.h"%g' {} +
```

files that added to the vendor folder:
```
vendor/github.com/ibm-messaging/mq-golang/v5/ibmmq/includes
 · cmqc.h
 · cmqcfc.h
 · cmqstrc.h
 · cmqxc.h
```
files that are patched:
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

#### :four: build the binary 

```bash
#inside the project folder
env GOOS=linux GOARCH=amd64 go build -mod=vendor ./...
```

The agent use an external config files locate in `conf` folder, there are secrets.json (ansible vault) with the credentials of the `elasticsearch` cluster and vars.json with the endpoint and a new parameter to select (in a future versions) the output of the agent, currently only is `elasticsearch`
 
---
## :large_blue_diamond: Install IBM MQ Server to deploy a full lab env and play

_first_ you have to delete `/opt/mqm` if you (like me) are recycling the machine, thats because the IBM MQ install process would like to use that folder, possible its can be change but I have no experience with this software so I'll try to keep all closer to the default. 

- Before you have to download the software from the IBM offical portal, you can register or use your IBM ID to get a Trial version.

- You can use the `apt` if you want, it is explained by IBM in his documentation too

```bash
dpkg -i ibmmq-runtime_9.2.0.0_amd64.deb
dpkg -i ibmmq-jre_9.2.0.0_amd64.deb
dpkg -i ibmmq-java_9.2.0.0_amd64.deb
dpkg -i ibmmq-gskit_9.2.0.0_amd64.deb
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

### Time to setup the IBM MQ to use the agent:

- Create the manager => `crtmqm -q gravity`
- Start the Manager => `strmqm gravity`
- Get into the console to create the queue => `runmqsc gravity`
- Create the queue => `define ql(gravity.cola01)`
- Check the queue => `dspmq`
- Add a message => `/opt/mqm/samp/bin/amqsput GRAVITY.COLA01` (2 empty enters to finish)

### Now is the moment to config the shell script 'mq_json.sh'
- The manager you will get metrics `queues="GRAVITY.*"`
- The path to the agent `exec /home/mqm/ibmmq-monitoring-agent $ARGS`

### Create the service on IBM MQ that will execute the shell script that execute the agent

```
DEFINE SERVICE(MQJSON)                         +
       CONTROL(QMGR)                           +
       SERVTYPE(SERVER)                        +
       STARTCMD('/home/mqm/mq_json.sh')        +
       STARTARG(+QMNAME+)                      +
       STOPCMD('/usr/bin/kill -9' )            +
       STOPARG(+MQ_SERVER_PID+)                +
       STDOUT('/var/mqm/errors/mq_json.out')   +
       STDERR('/var/mqm/errors/mq_json.err')   +
       DESCR('MQ exporter for JSON format')
```

Check if the service are correctly setup `DISPLAY SVSTATUS(MQJSON)`

Start the service `START SERVICE(MQJSON)`

Other useful commands:

```
STOP SERVICE(MQJSON)
DELETE SERVICE(MQJSON)
```

---
## If you want to run 'stan alone' the agent in a server with `IBM MQ`
```bash
./ibmmq-monitoring-agent -ibmmq.queueManager="gravity" -ibmmq.monitoredQueues="GRAVITY.*" -ibmmq.monitoredChannels="TO.*,SYSTEM.DEF.SVRCONN" -ibmmq.monitoredTopics="#" -ibmmq.monitoredSubscriptions="*" -ibmmq.interval="10s" -ibmmq.useStatus="true" -log.level="error"
```
remember that you need the `conf` folder in the same directory of the binary where are the secrets and config vars
---
---
# Always made with passion on golang !!!
<p align="center">
<img src="https://raw.githubusercontent.com/ocuil/assets/main/img/heart-hug.svg">
</p>
