# CustomResourceDefinition (CRD)
CRD are ways of creating new resource types in Kubernetes which follows same pattern as built-in types. 
We can combine custom resource with custom controller to encode the domain knowledge of specific application. 

When we create new CRD, API server creates a new RESTFUL resource path for each version you specify.

Reconstructing sample-controller using kubebuilder

Install kubebuilder
```
curl -L -O "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v1.0.8/kubebuilder_1.0.8_darwin_amd64.tar.gz"
tar -zxvf kubebuilder_1.0.8_darwin_amd64.tar.gz
mv kubebuilder_1.0.8_darwin_amd64 kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin
```


```
kubebuilder init --domain example.com
kubebuilder create api --group mysql --version v1alpha1 --kind Cluster
```

This creates 3 packages
cmd/... - contains manager main program. Manager is responsible for initializing shared dependencies and starting/stopping controllers. 

pkg/apis/... - contains api resource definitions. Edit "types.go" to implement API definitions.

pkg/controller/... - contains controller implementation. Edit "controller.go" to implement controller logic. 

config/... - contains yaml files for installing CRDs.


## Coding mysql operator
edit cfg/crds/mysql_cluster.yaml file 

edit pkg/api/mysql/v1alpha1/cluster_types.go 
