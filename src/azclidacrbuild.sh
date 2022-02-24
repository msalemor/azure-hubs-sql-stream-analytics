#!/bin/sh

# Build variables

version=1.5.1
acrname=alemoracr

# Build images

az acr build -r $acrname -f ./Sender.Dockerfile -t gosender:$version .
az acr build -r $acrname -f ./Receiver.Dockerfile -t goreceiver:$version .
az acr build -r $acrname -f ./EventsApi.Dockerfile -t goeventsapi:$version .
az acr build -r $acrname -f ./Monitor.Dockerfile -t gomonitor:$version .

# List image attributes 
#az acr repository list -n alemoracr

az acr repository show -n alemoracr --image gosender:$version
az acr repository show -n alemoracr --image goreceiver:$version
az acr repository show -n alemoracr --image goeventsapi:$version
az acr repository show -n alemoracr --image goeventsapi:$version
