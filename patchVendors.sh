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