#!/bin/bash

export MINIKUBE_WANTUPDATENOTIFICATION=false
export MINIKUBE_WANTREPORTERRORPROMPT=false
export CHANGE_MINIKUBE_NONE_USER=true
export KUBECONFIG=$HOME/.kube/config


# this for loop waits until kubectl can access the api server that Minikube has created
echo -n wait for minikube to start
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl version | grep "Server Version" > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n wait for local node to join
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl get no | grep " Ready " > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n wait for pods to start
for _ in {1..150}; do # timeout for 5 minutes
  if [ $(kubectl get --no-headers=true pods --all-namespaces=true 2>&1 | grep " Running " | grep "1/1" | wc -l) -eq "10" ]; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done
