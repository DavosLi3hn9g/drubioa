#!/bin/bash

PROC_NAME=iqiar
for i in  `ps -ef | grep -w $PROC_NAME | grep -v grep | awk '{print $2}'`;do sudo kill -9 $i;done
echo "QiarAI is stopped.."