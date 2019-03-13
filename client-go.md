# Programming Kubernetes

## Setup
Install k8s.io/client-go in your $GOPATH. 
```
# at the time of this post recommended tag was v10.0.0
go get k8s.io/client-go/v10.0.0
```

Install k8s.io/api-mechinery
```
go get -u k8s.io/apimachinery/v10.0.0
```

Dependency management
If you want to use dependencies from $GOPATH, then run 'godep restore' and remove vendor folder 
```
godep restore ./v10.0.0
rm -rf ./vendor
```

