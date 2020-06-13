#!/bin/bash

echo "Builds app"
go build -o live-omxctrl.bin

cd ./deploy

echo "build the zip package"
./deploy.bin -target pi4 -outdir ~/app/live-omxctrl/zips/
cd ~/app/live-omxctrl/

echo "update the service"
./update-service.sh

echo "Ready to fly"