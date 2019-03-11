# Operator Framework

Follow the pre-requisits and operator setup steps as in 
https://github.com/operator-framework/operator-sdk

mkdir -p $GOPATH/src/github.com/operators
cd $GOPATH/src/github.com/operator

operator-sdk new mysql-operator

operator-sdk add api --api-version=mysql.hashfab.io/v1alpha1 --kind=Mysql

operator-sdk add controller --api-version=mysql.hashfab.io/v1alpha1 --kind=Mysql

kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user deepakgb79@gmail.com

k create -f deploy/crds/mysql_v1alpha1_mysql_crd.yaml

operator-sdk build dgunjetti/mysql-operator

docker push dgunjetti/mysql-operator:latest

sed -i "" 's|REPLACE_IMAGE|dgunjetti/mysql-operator|g' deploy/operator.yaml

kubectl create -f deploy/service_account.yaml

kubectl create -f deploy/role.yaml

kubectl create -f deploy/role_binding.yaml

kubectl create -f deploy/operator.yaml

k create -f deploy/crds/mysql_v1alpha1_mysql_cr.yaml

