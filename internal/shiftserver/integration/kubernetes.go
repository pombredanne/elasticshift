/*
Copyright 2017 The Elasticshift Authors.
*/
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/conspico/elasticshift/api/types"
	itypes "gitlab.com/conspico/elasticshift/internal/shiftserver/integration/types"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type kubernetesClient struct {
	opts   *ConnectOptions
	Kube   *kubernetes.Clientset
	logger logrus.Logger
}

type ConnectOptions struct {
	Host                  string
	ServerCertificate     string
	Token                 string
	Namespace             string
	InsecureSkipTLSVerify bool

	Config    []byte
	UseConfig bool

	Storage types.Storage
}

var (
	//DefaultNamespace = "elasticshift"
	DefaultNamespace = "default"
	// DefaultNamespace = "shiftmk"
	DefaultContext = "elasticshift"

	KW_CREATEDBY = "createdby"
	KW_SHIFTID   = "shiftid"
	KW_BUILDID   = "buildid"
)

//type PersistentVolumeInfo struct {

//}

//type Client interface {
//	CreateContainer(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error)
//	CreateContainerWithVolume(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error)
//	CreatePersistentVolume(opts *CreatePersistentVolumeOption) (interface{}, error)

//	//DeleteContainer(info *ContainerInfo) (ContainerInfo, error)
//	//GetContainerStatus(opts *ContainerInfo) string
//}

func ConnectKubernetes(logger logrus.Logger, opts *ConnectOptions) (ContainerEngineInterface, error) {

	if opts.Namespace == "" {
		opts.Namespace = DefaultNamespace
	}

	kcli := &kubernetesClient{
		opts: opts,
	}

	var config *clientcmdapi.Config
	var err error

	if opts.UseConfig {

		// kube config file
		config, err = clientcmd.Load(opts.Config)
	} else {

		// use host, cert, user and token
		config = clientcmdapi.NewConfig()
		config.Clusters[DefaultContext] = &clientcmdapi.Cluster{
			Server: opts.Host,
			CertificateAuthorityData: []byte(opts.ServerCertificate),
		}

		config.AuthInfos[DefaultContext] = &clientcmdapi.AuthInfo{
			Token: opts.Token,
		}
		config.Contexts[DefaultContext] = &clientcmdapi.Context{
			Cluster:   DefaultContext,
			AuthInfo:  DefaultContext,
			Namespace: opts.Namespace,
		}
		config.CurrentContext = DefaultContext
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse the kube config: %v")
	}

	overrides := &clientcmd.ConfigOverrides{}
	if opts.InsecureSkipTLSVerify {
		overrides.ClusterInfo = clientcmdapi.Cluster{
			InsecureSkipTLSVerify: true,
		}
	}

	clientBuilder := clientcmd.NewDefaultClientConfig(*config, overrides)

	clientConfig, err := clientBuilder.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to kubernetes : %v", err)
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	kcli.Kube = clientset
	return kcli, nil
}

func (c *kubernetesClient) CreateContainer(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error) {

	envs := []apiv1.EnvVar{}
	for _, env := range opts.Environment {
		envs = append(envs, apiv1.EnvVar{Name: env.Key, Value: env.Value})
	}

	volumeMounts := []apiv1.VolumeMount{}
	volumes := []apiv1.Volume{}
	for _, vol := range opts.VolumeMounts {
		volumeMounts = append(volumeMounts, apiv1.VolumeMount{Name: vol.Name, MountPath: vol.MountPath})

		if c.opts.Storage.Kind == 4 {

			// vol := apiv1.Volume{
			// 	Name: vol.Name,
			// 	VolumeSource: apiv1.VolumeSource{
			// 		NFS: &apiv1.NFSVolumeSource{
			// 			Server:   c.opts.Storage.NFS.Server,
			// 			Path:     c.opts.Storage.NFS.Path,
			// 			ReadOnly: c.opts.Storage.NFS.ReadOnly,
			// 		},
			// 	},
			// }

			vol := apiv1.Volume{
				Name: vol.Name,
				VolumeSource: apiv1.VolumeSource{
					HostPath: &apiv1.HostPathVolumeSource{
						Path: "/opt/elasticshift",
					},
				},
			}
			volumes = append(volumes, vol)
		}
		// volumes = append(volumes, apiv1.Volume{Name: vol.Name, VolumeSource: apiv1.VolumeSource{PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{ClaimName: vol.Claim}}})
	}

	deploymentsClient := c.Kube.AppsV1().Deployments(c.opts.Namespace)

	shiftId := uuid.NewV4()
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: opts.BuildID,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					KW_BUILDID: opts.BuildID,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						KW_CREATEDBY: DefaultContext,
						KW_SHIFTID:   shiftId.String(),
						KW_BUILDID:   opts.BuildID,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:         opts.BuildID,
							Image:        opts.Image,
							Command:      []string{opts.Command},
							Env:          envs,
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}

	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		return nil, fmt.Errorf("Error in creating container : %v", err)
	}
	c.logger.Infof("Created deployment %q.\n", result.GetObjectMeta().GetName())

	watch, err := deploymentsClient.Watch(v1.ListOptions{LabelSelector: "createdby=" + DefaultContext})
	lo := &v1.ListOptions{LabelSelector: KW_SHIFTID + "=" + shiftId.String()}
	fmt.Printf("List options : %v", lo)
	watch, err = c.Kube.CoreV1().Pods(c.opts.Namespace).Watch(*lo)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	for {
		select {
		case res := <-watch.ResultChan():
			z, err := json.Marshal(res)

			if err != nil {
				fmt.Errorf("%v", err)
			}
			var out bytes.Buffer
			json.Indent(&out, z, "=", "\t")
			out.WriteTo(os.Stdout)
			//fmt.Println(res.Type, string(z))

			if res.Type == "Modified" {
				// Stop when the status changed to modified, in real need to check the status Running and then
				// this should be stopped
				watch.Stop()
			}
		}
	}
	watch.Stop()

	//Construct ContainerInfo
	md := result.GetObjectMeta()
	cinfo := &itypes.ContainerInfo{
		Name: md.GetName(),
		// CreationTimestamp: md.GetCreationTimestamp().String(),
		Status:       result.Status.String(),
		Image:        opts.Image,
		ImageVersion: opts.ImageVersion,
		ClusterName:  md.GetClusterName(),
		UID:          string(md.GetUID()),
		ShiftID:      shiftId.String(),
		Namespace:    md.GetNamespace(),
	}

	return cinfo, nil
}

func (c *kubernetesClient) DeleteContainer(id string) error {

	deletePolicy := metav1.DeletePropagationForeground
	deploymentsClient := c.Kube.AppsV1().Deployments(c.opts.Namespace)
	err := deploymentsClient.Delete(id, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	return err
}

//watch, err := deploymentsClient.Watch(v1.ListOptions{LabelSelector: "createdby=elasticshift"})
//lo := &v1.ListOptions{LabelSelector: "esuuid=" + uid.String()}
//fmt.Printf("List options : %v", lo)
//watch, err := c.Kube.CoreV1().Pods(namespace).Watch(*lo)
//if err != nil {
//	fmt.Errorf("%v", err)
//}
//for {
//	select {
//	//case res := <-watch.ResultChan():
//	case res := <-watch.ResultChan():
//		//e := res.(watch.Event)
//		z, err := json.Marshal(res)

//		if err != nil {
//			fmt.Errorf("%v", err)
//		}
//		var out bytes.Buffer
//		json.Indent(&out, z, "=", "\t")
//		out.WriteTo(os.Stdout)
//		//fmt.Println(res.Type, string(z))

//		//if res.Type == "MODIFIED" {
//		//Stop when the status changed to modified, in real need to check the status Running and then
//		//this should be stopped
//		//watch.Stop()
//		//}

//	}

//}
//watch.Stop()

func (c *kubernetesClient) CreateContainerWithVolume(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error) {

	deploymentsClient := c.Kube.AppsV1beta1().Deployments(c.opts.Namespace)

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

	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		fmt.Errorf("Error in creating container : %v", err)
	}
	c.logger.Infof("Created deployment %q.\n", result.GetObjectMeta().GetName())

	//Construct ContainerInfo
	md := result.GetObjectMeta()
	cinfo := &itypes.ContainerInfo{
		Name:              md.GetName(),
		CreationTimestamp: md.GetCreationTimestamp().String(),
		//StoppedAt:         nil,
		Status:       result.Status.String(),
		Image:        opts.Image,
		ImageVersion: opts.ImageVersion,
		ClusterName:  md.GetClusterName(),
		UID:          string(md.GetUID()),
		Namespace:    md.GetNamespace(),
	}

	//watch, err := deploymentsClient.Watch(v1.ListOptions{LabelSelector: "createdby=elasticshift"})
	lo := &v1.ListOptions{LabelSelector: "esuuid=" + uid.String()}
	fmt.Printf("List options : %v", lo)
	watch, err := c.Kube.CoreV1().Pods(apiv1.NamespaceDefault).Watch(*lo)
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

func (c *kubernetesClient) PersistentVolumeClaim(opts *itypes.PersistentVolumeClaimOptions) (interface{}, error) {

	q, err := resource.ParseQuantity(opts.Capacity)
	if err != nil {
		return nil, err
	}

	pvc := c.Kube.CoreV1().PersistentVolumeClaims(c.opts.Namespace)
	persistentVolumeClaim, err := pvc.Create(&apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: opts.Name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteMany},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceRequestsStorage: q,
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}
	return persistentVolumeClaim, nil
}

func (c *kubernetesClient) CreatePersistentVolume(opts *itypes.CreatePersistentVolumeOptions) (*itypes.PersistentVolumeInfo, error) {

	q, err := resource.ParseQuantity(opts.Capacity)
	if err != nil {
		panic(err) // handle it
	}

	pv := c.Kube.CoreV1().PersistentVolumes()
	persistentVolume, err := pv.Create(&apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: opts.Name,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		Spec: apiv1.PersistentVolumeSpec{
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				NFS: &apiv1.NFSVolumeSource{
					Server:   opts.Server,
					Path:     opts.Path,
					ReadOnly: false,
				},
			},
			AccessModes:                   []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteMany},
			MountOptions:                  opts.MountOptions,
			Capacity:                      map[apiv1.ResourceName]resource.Quantity{apiv1.ResourceStorage: q},
			PersistentVolumeReclaimPolicy: apiv1.PersistentVolumeReclaimRetain,
		},
	})

	if err != nil {
		return nil, err
	}

	vi := &itypes.PersistentVolumeInfo{
		Name: persistentVolume.Name,
	}

	return vi, nil
}

//func (c *kubernetesClient) GetContainerStatus(opts *CreateContainerOptions) string {

//return ""
//}

func int32Ptr(i int32) *int32 { return &i }
