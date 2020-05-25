# k8s-patch
Kubernetes Patch / Apply in Go

This project is to simulate the `kubectl apply` when you are creating a operator and need to apply the difference between the current and modified resource

Creating the resource of first time and setting the metadata `metadata.SetMetadata(obj)`

```go
import "github.com/oassuncao/k8s-patch/metadata"
...

if err := metadata.SetMetadata(obj); err != nil {
  log.Error(err, "Error on setting annotation data")
  return reconcile.Result{}, err
}
  
err := r.client.Create(context.TODO(), obj)
if err != nil {
  log.Error(err, "Error on creating")
  return reconcile.Result{}, err
}
```

Checking if has any change
```go
 import "github.com/oassuncao/k8s-patch/compare"
 ...
 
 
found := &corev1.Service{}
err := r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: instance.Namespace}, found)
if err == nil {
		equal, err := compare.DeepEqualPatch(found, service);
		if err != nil {
			log.Error(err, "Error on DeepEqual")
			return reconcile.Result{}, err
		}

		if equal {
			log.Info("Service already exists")
			return reconcile.Result{}, nil
		}

    log.Info("The Service is changed")
}
```

Creating the Patch to be apply
```go
import "github.com/oassuncao/k8s-patch/metadata"
...

patch, err := metadata.GeneratePatch(found, service)
if err != nil {
  log.Error(err, "Error on Generate Patch")
  return reconcile.Result{}, err
}

err := r.client.Patch(context.TODO(), service, client.RawPatch(patch.PatchType, patch.Data))
if err != nil {
  log.Error(err, "Error on patching")
  return reconcile.Result{}, err
}
```
