#!/bin/sh
case $1 in
        start)
                nohup ./im_go 2>&1 >> info.log 2>&1 /dev/null &
                echo "服务已启动..."
                sleep 1
        ;;
        stop)
                killall im_go
                echo "服务已停止..."
                sleep 1
        ;;
        restart)
                killall im_go
                sleep 1
                nohup ./im_go 2>&1 >> info.log 2>&1 /dev/null &
                echo "服务已重启..."
                sleep 1
        ;;
        *)
                echo "$0 {start|stop|restart}"
                exit 4
        ;;
esac
