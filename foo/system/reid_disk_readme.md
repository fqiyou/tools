#### raid disk 说明
>根据华为云的shell脚本改写的。由于目前我们使用监控为 xx_export+prometheus+grafana,故开发raid_disk_export

##### raid_disk_export地址

```cgo
https://github.com/fqiyou/tools/blob/master/foo/monitor/host/raid_disk/main.go
```
##### 华为云的shell地址
```
https://support-intl.huaweicloud.com/zh-cn/usermanual-ces/zh-cn_topic_0101158226.html
```

##### 华为云的shell
```
#!/bin/bash
##################################################
## description: check raid status script(mdadm tool)
## create time: 2018/3/16
#################################################
debug=false
if [ -n $1 ] && [ "$1" = "-d" ];then
    debug=true
elif [ "$1" = "-x" ]; then
    #statements
        set -x
fi

main()
{
    get_all_md
    outStr='{"data":[ '
    declare -A mdResult=()
    for md in $all_md
    do
        get_md_device_statistic $md
        get_md_status $md
    done

    for key in ${!mdResult[@]}
    do
        if [ -z ${mdResult[$key]} ];then
            mdResult[$key]=0
        fi
        metric_prefix=$(echo $key | cut -d "@" -f1)
        metric_name=$(echo $key | cut -d "@" -f2)
        outStr=$outStr'{"metric_name" : "'$metric_name'","metric_value" : '${mdResult[$key]}',"metric_prefix" : "'$metric_prefix'"},'
    done

    seconds=`date +%s`000
    ## 指标去掉最后一个逗号
    outStr=${outStr%?}' ], ''"collect_time":'$seconds'}'

    echo $outStr
    return 0
}

get_md_status()
{
    mdResult[$1@status_device]=1
    # check md
    mdStatus=$(cat /proc/mdstat |grep $1|grep -P '\sactive\s'|tr -d ' ')
    if [ -z $mdStatus ];then
        mdResult[$1@status_device]=0
        return
    fi
    degraded=$(mdadm -D /dev/$1|grep State|grep degraded|tr -d ":"|tr -d " ")
    if [ -n "$degraded" ];then
        mdResult[$1@status_device]=0
        return
    fi
    # check device
    devices=$(cat /proc/mdstat|grep $1' '|grep raid|awk -F'raid[0-9]+' '{print $2}'|sed 's/\[[0-9]*\]//g'|sed 's/(.)//g'|xargs)
    for dd in $devices
    do
        deviceResult=$(mdadm -D /dev/$1|awk 'BEGIN{FS="\n";RS="\n\n"}{pp[NR]=$0; if ($0 ~ /RaidDevice/) {line=NR} }END{for(i=line;i<=NR;i++) print pp[i]}'|grep $dd|grep active|tr -d ' ')
        if [[ -z $deviceResult ]]; then
            #device adnormal
            mdRsult[$1@status_device]=0
            break
        fi
    done
}

get_md_device_statistic()
{
    mdResult[$1"@active_device"]=$(mdadm -D /dev/$1 |grep -i "Active Devices" |awk -F ":" '{print $2}' |tr -d " ")
    mdResult[$1"@working_device"]=$(mdadm -D /dev/$1 |grep -i "Working Devices" |awk -F ":" '{print $2}' |tr -d " ")
    mdResult[$1"@failed_device"]=$(mdadm -D /dev/$1 |grep -i "Failed Devices" |awk -F ":" '{print $2}' |tr -d " ")
    mdResult[$1"@spare_device"]=$(mdadm -D /dev/$1 |grep -i "Spare Devices" |awk -F ":" '{print $2}' |tr -d " ")
}

get_all_md()
{
    #
    all_md=$(cat /proc/mdstat |grep md|awk -F ":" '{print $1}')
    if [ 0 -ne $? ] || [ -z "$all_md" ]
    then
        if $debug;then
           echo "[ERROR] please check is there any raid!"
        fi
        exit 0
    fi
    if $debug;then
        for md in $all_md
        do
           echo "md device:"
           echo $md
        done
    fi
    return 0
}

main
exit $?


```