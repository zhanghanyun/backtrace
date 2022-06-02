#!/bin/bash

arch=$(uname -m)
if [ "$arch" = "x86_64" ]; then
  wget -q -O backtrace.tar.gz  https://github.com/zhanghanyun/backtrace/releases/latest/download/backtrace-linux-amd64.tar.gz
else
  wget -q -O backtrace.tar.gz  https://github.com/zhanghanyun/backtrace/releases/latest/download/backtrace-linux-arm64.tar.gz
fi

tar -xf backtrace.tar.gz
rm backtrace.tar.gz
mv backtrace /usr/bin/
backtrace