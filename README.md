## Create a Controller using Kubebuilder

Controller enable us to create a resource with `Kind: SampleVolume` which in turn create a PVC.

1. Reconcile
2. SetControllerReference for OwnerReference

## Steps

1. Create scaffolding

```bash
kubebuilder init --domain operator.yogeshsharma.me --repo demo-volume
```

2. Create API/Custom Resources

```bash
kubebuilder create api --group samplevolumne --version v1 --kind SampleVolume
```


3. Update `api/v1/samplevolume_types.go` and execute make commands

```go
type SampleVolumeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SampleVolume. Edit samplevolume_types.go to remove/update
	Name string `json:"name,omitempty"`
	Size int 	`json:"size,omitempty"`
}

// SampleVolumeStatus defines the observed state of SampleVolume
type SampleVolumeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status,omitempty"`
}
```

`make generate` --> This will generate a CRD in `config/patches/webhook_in_samplevolumes.yaml`

`make manifests` --> Create a manifests in `config/crd/bases` which we can load directly in k8s.


4. Controller Code

It is written in `controllers/samplevolume_controller.go` under `Reconcile` for Reconcile Loop.

After writing the logic above, `make install` which will install CRD into cluster.

Verify using `k8 get crd`


```bash
$ k8 get crd
NAME                                                   CREATED AT
samplevolumes.samplevolumne.operator.yogeshsharma.me   2022-03-24T06:02:25Z
```

5. A sample resource definition is present in `config/samples/samplevolumne_v1_samplevolume.yaml`


```bash
$ k8 apply -f config/samples/samplevolumne_v1_samplevolume.yaml 
samplevolume.samplevolumne.operator.yogeshsharma.me/samplevolume-sample created
```

Check status

```bash
$ k8 get samplevolumes.samplevolumne.operator.yogeshsharma.me
NAME                  AGE
samplevolume-sample   56s
```

6. `make run` to check log and run in foreground.

7. To run in cluster

```bash
make docker-build docker-push IMG=yks0000/sample-volume-k8s-operator:0.1
make deploy IMG=yks0000/sample-volume-k8s-operator:0.1
```


8. For cluster Installation we need to have PVC RBAC so that service account can watch/update/create/delete PVC when SampleVolume event triggered.

Files:

```bash
$ ll config/rbac/pvc_role*
-rw-r--r--  1 yosharma  staff   662B Mar 25 11:03 config/rbac/pvc_role.yaml
-rw-r--r--  1 yosharma  staff   280B Mar 25 11:04 config/rbac/pvc_role_binding.yaml
```

9. `SetControllerReference` for OwnerReference.

This is used for garbage collection of the controlled object and for reconciling the owner object on changes to controlled (with a Watch + EnqueueRequestForOwner).

```go
logger.Info("Setting up pvc controller reference")
if err = controllerutil.SetControllerReference(volume, pvc, r.Scheme); err != nil {
    logger.Error(err, "Failed to set pvc controller reference")
    return err
	}
```