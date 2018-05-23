#!/bin/bash
#===============================================================
export HEADER='\033[95m'
export OKBLUE='\033[94m'
export OKGREEN='\033[92m'
export WARNING='\033[93m'
export FAIL='\033[91m'
export ENDC='\033[0m'
export BOLD='\033[1m'
export UNDERLINE='\033[4m'
#===============================================================

V=$(date "+%Y%m%d_%H%M%S")
PROJECT="neuron-191011"
NAMESPACE=staging-exchange
BACKEND_IMAGE="$NAMESPACE-crypto-exchange"

gcloud auth activate-service-account --key-file ./credentials/deploy.cred.json
gcloud container clusters get-credentials neuron-cluster-1 --zone us-west1-a --project neuron-191011

if [ $1 = "staging" ]
then
    cp -a ./credentials/staging.cred.json ./credentials/cred.json
fi

if [ $1 = "production" ]
then
    cp -a ./credentials/production.cred.json ./credentials/cred.json
fi

buildNumber=$V
docker build \
    -t gcr.io/$PROJECT/$BACKEND_IMAGE:$buildNumber .
docker tag gcr.io/$PROJECT/$BACKEND_IMAGE:$buildNumber gcr.io/$PROJECT/$BACKEND_IMAGE:$buildNumber

gcloud docker -- push gcr.io/$PROJECT/$BACKEND_IMAGE:$buildNumber

result=$(echo $?)
if [ $result != 0 ] ; then
    echo "$FAIL failed gcloud docker -- push gcr.io/$PROJECT/$BACKEND_IMAGE:buildNumber $V $ENDC";
    exit;
else
    echo "$OKGREEN gcloud docker -- push gcr.io/$PROJECT/$BACKEND_IMAGE:buildNumber $V $ENDC"
fi

kubectl --namespace=$NAMESPACE set image deployment/backend backend=gcr.io/$PROJECT/$BACKEND_IMAGE:$buildNumber

result=$(echo $?)
if [ $result != 0 ] ; then
    echo "$FAIL failed kubectl --namespace=$NAME_SPACE set image deployment/backend backend=gcr.io/$PROJECT/$BACKEND_IMAGE:$buildNumber $ENDC";
    exit;
else
    echo "$OKGREEN DEPLOY SUCESSFULL gcr.io/$PROJECT/$BACKEND_IMAGE:$buildNumber $ENDC"
fi