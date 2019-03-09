# Operator Framework

Follow the pre-requisits and operator setup steps as in 
https://github.com/operator-framework/operator-sdk

mkdir -p $GOPATH/src/github.com/operators
cd $GOPATH/src/github.com/operator

operator-sdk new simple-operator --skip-git-init

kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user deepakgb79@gmail.com


