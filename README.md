## Create a PVC with CRD.


1. Crate scaffolding

kubebuilder init --domain operator.yogeshsharma.me --repo demo-volume

2. Create API/Custom Resources

kubebuilder create api --group samplevolumne --version v1 --kind SampleVolume

3. Update api/v1/samplevolume_types.go and execute make commands

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

make docker-build docker-push IMG=yks0000/sample-volume-k8s-operator:0.1
make deploy IMG=yks0000/sample-volume-k8s-operator:0.1
