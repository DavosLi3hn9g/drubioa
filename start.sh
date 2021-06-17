#!/bin/bash

PROC_NAME=iqiar
ProcNumber=$(ps -ef | grep -w $PROC_NAME | grep -v grep | wc -l)
if [ $ProcNumber -le 0 ]; then
  sudo chmod 777 ./stop.sh
  sudo chmod +x ./iqiar
  sudo chmod +x ./reload
  sudo chmod +rw ./data/*
  sudo nohup ./iqiar $* > nohup.out 2>&1 &
  echo "QiarAI is started.."
else
  echo "QiarAI is running.."
fi
