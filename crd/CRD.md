# CustomResourceDefinition (CRD)
CRD are ways of creating new resource types in Kubernetes which follows same pattern as built-in types. 
We can combine custom resource with custom controller to encode the domain knowledge of specific application. 

When we create new CRD, API server creates a new RESTFUL resource path for each version you specify.

sample-controller on Kubernetes Github

```
cat << EOF > foo-crd.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: foos.samplecontroller.k8s.io
spec:
  group: samplecontroller.k8s.io
  version: v1alpha1
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
EOF
```

```
kubectl create -f foo-crd.yaml
```



