# Kubernetes Custom Controllers

Controllers are control loops that watches the shared state of cluster through API server and makes changes to move the current state to desired state. 
Ex: Deployment controller, Replicaset controller, node controller etc..

We can create our own controllers for core and custom resources.

## Controller componentes

### SharedInformer
Informer is a shared data cache. Informer deliver events to listeners when the data they are interested in changes. We create reference to Informer using factory method.
```
podInformer = InformerFactory.Core().V1().Pods()
```

### Resource Event Handlers
Controllers register interest in a specific object using EventHandlers. Controller can handle create/update/delete events.
We register callback functions which will be called by Informers to deliver events to controller. 
```
podInformer.informer().AddEventHandler(
  cache.ResourceHandler.Funcs{
    AddFunc: func(obj Interface{}) {}
    UpdateFunc: func(old, cur Interface{}) {}
    DeleteFunc: func(obj Interface{}) {}
  }
)
```

### Workqueue
Resource Event Handler callback functions obtain an object key from the events, enqueue that key to a work queue for further processing. Object key is combination of namespace and name of resource.

### Workers
Workers are go routines where we processes items in work queue, we can use index reference or listing wrapper to retrieve object from object key. It contains the sync handler where we run business logic of the controller.

```
func(c *ctrl) worker() {
  for c.ProcessNextItem() {}
}

func (c *ctrl) processNextItem() {
  item := c.queue.Get()
  err := c.syncHandler()
  c.handlerErr(err, item)
}
```

