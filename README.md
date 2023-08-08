# Namespace Label Propagator

This policy ensures that all workloads deployed in the cluster contain specific
labels defined in the namespace where the workload will be executed.

When a resource is created or updated, this policy copies a list of user-specified
labels from the namespace to the object. Labels defined on the namespace take
precedence over labels already defined inside the resource.

This policy is able to set the labels for the following resource kinds: `Pod`,
`ReplicationController`, `Deployment`, `ReplicaSet`, `StatefulSet`, `DaemonSet`,
`Job` and `CronJob`.

When dealing with Kubernetes resources that generate pods, the policy ensures the
special labels are propagated also to them.

## Settings

This policy has a single setting called `propagatedLabels`, which is a list of
strings representing the labels from the namespace definition that should be
propagated to the workloads deployed in the namespace.

For example, the policy configuration would look like this:

```yaml
propagatedLabels:
- cost-center
- field.cattle.io/projectId
```

In this scenario, when a resource is created, the policy ensures that the
`cost-center` and `field.cattle.io/projectId` labels are copied from the
namespace object to the resource itself.

Label propagation only occurs if the desired labels are already set on the namespace.
If a label is not defined in the namespace, it will not be propagated to the workloads
