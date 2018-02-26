#!/bin/bash

export MINIKUBE_WANTUPDATENOTIFICATION=false
export MINIKUBE_WANTREPORTERRORPROMPT=false
export MINIKUBE_HOME=$HOME
export CHANGE_MINIKUBE_NONE_USER=true
export KUBECONFIG=$HOME/.kube/config


# this for loop waits until kubectl can access the api server that Minikube has created
echo -n wait for minikube to start
for i in {1..150}; do # timeout for 5 minutes
  kubectl version | grep "Server Version" > /dev/null 2>&1
  if [ $? -eq 0 ]; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n wait for local node to join
for i in {1..150}; do # timeout for 5 minutes
  kubectl get no | grep " Ready " > /dev/null 2>&1
  if [ $? -eq 0 ]; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n wait for dns to start
for i in {1..150}; do # timeout for 5 minutes
  kubectl get --no-headers=true pods -n kube-system -l k8s-app=kube-dns | grep " Running " | grep "3/3" > /dev/null 2>&1
  if [ $? -eq 0 ]; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

# kubectl commands are now able to interact with Minikube cluster

# workaround https://github.com/kubernetes/minikube/issues/1947
echo -n getting name of kubedns pod
for i in {1..150}; do # timeout for 5 minutes
  KUBEDNS_POD=$(kubectl get --no-headers=true pods -n kube-system -l k8s-app=kube-dns -o custom-columns=:metadata.name)
  if [ $? -eq 0 ]; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n fixing kubedns upstream server
for i in {1..150}; do # timeout for 5 minutes
  kubectl exec -n kube-system $KUBEDNS_POD -c kubedns -- sh -c "echo nameserver 8.8.8.8 > /etc/resolv.conf"
  if [ $? -eq 0 ]; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done
