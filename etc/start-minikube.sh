#!/bin/sh

export MINIKUBE_WANTUPDATENOTIFICATION=false
export MINIKUBE_WANTREPORTERRORPROMPT=false
export MINIKUBE_HOME=$HOME
export CHANGE_MINIKUBE_NONE_USER=true
mkdir -p $HOME/.kube
touch $HOME/.kube/config

export KUBECONFIG=$HOME/.kube/config
sudo -E minikube start --vm-driver=none

# this for loop waits until kubectl can access the api server that Minikube has created
for i in {1..150}; do # timeout for 5 minutes
   kubectl get po &> /dev/null
   if [ $? -ne 1 ]; then
      break
  fi
  sleep 2
done

# kubectl commands are now able to interact with Minikube cluster

# workaround https://github.com/kubernetes/minikube/issues/1947
echo fixing kubedns in minikube
KUBEDNS_POD=$(kubectl get --no-headers=true pods -n kube-system -l k8s-app=kube-dns -o custom-columns=:metadata.name)
kubectl exec -n kube-system $KUBEDNS_POD -c kubedns -- sh -c "echo nameserver 8.8.8.8 > /etc/resolv.conf"
