/*
Copyright 2017 The Elasticshift Authors.
*/
package kubernetes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type kubernetesClient struct {
	Kube *kubernetes.Clientset
}

//type PersistentVolumeInfo struct {

//}

type KubernetesClient interface {
	CreateContainer(opts *CreateContainerOptions) (*ContainerInfo, error)
	CreateContainerWithVolume(opts *CreateContainerOptions) (*ContainerInfo, error)
	CreatePersistentVolume(opts *CreatePersistentVolumeOption) (interface{}, error)

	//DeleteContainer(info *ContainerInfo) (ContainerInfo, error)
	//GetContainerStatus(opts *ContainerInfo) string
}

func NewKubernetesClient(o *KubernetesClientOptions) (KubernetesClient, error) {

	if o.KubeConfigFile == "" {
		errors.New("Kubernetes config file required to proceed further")
	}

	f, err := os.Open(o.KubeConfigFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to open Kubernetes config file %s : %v", o.KubeConfigFile, err)
	}
	defer f.Close()

	//if home := homedir.HomeDir(); home != "" {
	//kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//} else {
	//kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	//}
	//flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", o.KubeConfigFile)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &kubernetesClient{Kube: clientset}, nil

}

func (kc *kubernetesClient) CreateContainer(opts *CreateContainerOptions) (*ContainerInfo, error) {
	//fmt.Printf("kube config : %v", kc.Kube)
	r, _ := json.Marshal(kc.Kube)
	fmt.Printf("The Marshal string %v", string(r))

	deploymentsClient := kc.Kube.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	uid := uuid.NewV4()
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "es-build-" + opts.Image + opts.ImageVersion,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"createdby": "elasticshift",
						"esuuid":    uid.String(),
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  opts.Image,
							Image: opts.Image + ":" + opts.ImageVersion,
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating container...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		fmt.Errorf("Error in creating container : %v", err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	//Construct ContainerInfo
	md := result.GetObjectMeta()
	cinfo := &ContainerInfo{
		Name:              md.GetName(),
		CreationTimestamp: md.GetCreationTimestamp().String(),
		//StoppedAt:         nil,
		Status:       result.Status.String(),
		Image:        opts.Image,
		ImageVersion: opts.ImageVersion,
		ClusterName:  md.GetClusterName(),
		Uid:          string(md.GetUID()),
		Namespace:    md.GetNamespace(),
	}
	//podTemplateHash := result.GetLabels()["pod-template-hash"]
	//podTemplateHash := result.GetLabels()
	//g, err := result.Marshal()
	//if err != nil {
	//panic(err)
	//}
	status := result
	//status := result.GetInitializers().Result.Status(),ci
	s, _ := json.Marshal(status)
	//var out bytes.Buffer
	//json.Indent(&out, status, "=", "\t")
	//out.WriteTo(os.Stdout)
	fmt.Println("STATUS : %v", string(s))

	//fmt.Println("result.Marshal() : %v", string(g))
	//fmt.Println("PodTemplateHash : %v", podTemplateHash   )
	//v1.ListOptions(LabelSelector )

	//watch, err := deploymentsClient.Watch(v1.ListOptions{LabelSelector: "createdby=elasticshift"})
	lo := &v1.ListOptions{LabelSelector: "esuuid=" + uid.String()}
	fmt.Printf("List options : %v", lo)
	watch, err := kc.Kube.CoreV1().Pods(apiv1.NamespaceDefault).Watch(*lo)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	for {
		select {
		//case res := <-watch.ResultChan():
		case res := <-watch.ResultChan():
			//e := res.(watch.Event)
			z, err := json.Marshal(res)

			if err != nil {
				fmt.Errorf("%v", err)
			}
			var out bytes.Buffer
			json.Indent(&out, z, "=", "\t")
			out.WriteTo(os.Stdout)
			//fmt.Println(res.Type, string(z))

			//if res.Type == "MODIFIED" {
			//Stop when the status changed to modified, in real need to check the status Running and then
			//this should be stopped
			//watch.Stop()
			//}

		}

	}
	//watch.Stop()
	return cinfo, nil
}

func (kc *kubernetesClient) CreateContainerWithVolume(opts *CreateContainerOptions) (*ContainerInfo, error) {

	//fmt.Printf("kube config : %v", kc.Kube)
	r, _ := json.Marshal(kc.Kube)
	fmt.Printf("The Marshal string %v", string(r))

	deploymentsClient := kc.Kube.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	uid := uuid.NewV4()
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "es-build-" + opts.Image + opts.ImageVersion,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"createdby": "elasticshift",
						"esuuid":    uid.String(),
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  opts.Image,
							Image: opts.Image + ":" + opts.ImageVersion,
							// EnvVar: opts.Environment, TODO
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating container...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		fmt.Errorf("Error in creating container : %v", err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	//Construct ContainerInfo
	md := result.GetObjectMeta()
	cinfo := &ContainerInfo{
		Name:              md.GetName(),
		CreationTimestamp: md.GetCreationTimestamp().String(),
		//StoppedAt:         nil,
		Status:       result.Status.String(),
		Image:        opts.Image,
		ImageVersion: opts.ImageVersion,
		ClusterName:  md.GetClusterName(),
		Uid:          string(md.GetUID()),
		Namespace:    md.GetNamespace(),
	}
	//podTemplateHash := result.GetLabels()["pod-template-hash"]
	//podTemplateHash := result.GetLabels()
	//g, err := result.Marshal()
	//if err != nil {
	//panic(err)
	//}
	status := result
	//status := result.GetInitializers().Result.Status(),ci
	s, _ := json.Marshal(status)
	//var out bytes.Buffer
	//json.Indent(&out, status, "=", "\t")
	//out.WriteTo(os.Stdout)
	fmt.Println("STATUS : %v", string(s))

	//fmt.Println("result.Marshal() : %v", string(g))
	//fmt.Println("PodTemplateHash : %v", podTemplateHash   )
	//v1.ListOptions(LabelSelector )

	//watch, err := deploymentsClient.Watch(v1.ListOptions{LabelSelector: "createdby=elasticshift"})
	lo := &v1.ListOptions{LabelSelector: "esuuid=" + uid.String()}
	fmt.Printf("List options : %v", lo)
	watch, err := kc.Kube.CoreV1().Pods(apiv1.NamespaceDefault).Watch(*lo)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	for {
		select {
		//case res := <-watch.ResultChan():
		case res := <-watch.ResultChan():
			//e := res.(watch.Event)
			z, err := json.Marshal(res)

			if err != nil {
				fmt.Errorf("%v", err)
			}
			var out bytes.Buffer
			json.Indent(&out, z, "=", "\t")
			out.WriteTo(os.Stdout)
			//fmt.Println(res.Type, string(z))

			//if res.Type == "MODIFIED" {
			//Stop when the status changed to modified, in real need to check the status Running and then
			//this should be stopped
			//watch.Stop()
			//}

		}

	}
	//watch.Stop()
	return cinfo, nil
}
func (kc *kubernetesClient) CreatePersistentVolume(opts *CreatePersistentVolumeOption) (interface{}, error) {
	pv := kc.Kube.CoreV1().PersistentVolumes()

	q, err := resource.ParseQuantity(opts.Capacity)

	if err != nil {
		panic(err) // handle it
	}

	persistentVolume, err := pv.Create(&apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: opts.Name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		Spec: apiv1.PersistentVolumeSpec{
			//Capacity :
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				NFS: &apiv1.NFSVolumeSource{
					Server: opts.Server,
					Path:   opts.Path,
				},
			},
			AccessModes:  []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteMany},
			MountOptions: opts.MountOptions,
			//Capacity:     map[apiv1.ResourceName]resource.Quantity{apiv1.ResourceStorage: *resource.NewQuantity(5, resource.BinarySI)},
			Capacity: map[apiv1.ResourceName]resource.Quantity{apiv1.ResourceStorage: q},
			//MountOptions: []string{"hard", "nfsvers=4.1"},
		},
		//Spec : apiv1.PersistentVolumeSpec{

		//}
	})

	if err != nil {
		panic(err)
	}

	return persistentVolume, nil
}

//func (kc *kubernetesClient) GetContainerStatus(opts *CreateContainerOptions) string {

//return ""
//}

func int32Ptr(i int32) *int32 { return &i }
