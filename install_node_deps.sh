#!/bin/sh

wget https://github.com/prometheus/node_exporter/releases/download/v1.3.1/node_exporter-1.3.1.linux-amd64.tar.gz
tar xvf node_exporter-1.3.1.linux-amd64.tar.gz
cd node_exporter-1.3.1.linux-amd64/
sudo cp node_exporter /usr/local/bin
cd ..
rm -rf node_exporter-1.3.1.linux-amd64
rm -rf node_exporter-1.3.1.linux-amd64.tar.gz

sudo useradd --no-create-home --shell /bin/false node_exporter
sudo chown node_exporter:node_exporter /usr/local/bin/node_exporter
sudo vim /etc/systemd/system/node_exporter.service