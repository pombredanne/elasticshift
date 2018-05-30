/*
Copyright 2017 The Elasticshift Authors.
*/
package integration

import (
	"testing"
)

func testNewKubernetesClient(t *testing.T) {

	// kc := getKubernetesClient()

	// fmt.Println("Kub Client", kc)
	//c, err := kc.CreateContainer(&CreateContainerOptions{Image: "nginx", ImageVersion: "1.13.6"})

	// if err != nil {
	// 	panic(err)
	// }

	// j, err := json.Marshal(c)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("%v", string(j))

}

func testCreateContainer2(t *testing.T) {

	// host := "https://192.168.99.100:8443"
	// cacert := `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQVRBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwdGFXNXAKYTNWaVpVTkJNQjRYRFRFNE1ERXdPVEV5TlRFeE5Wb1hEVEk0TURFd056RXlOVEV4TlZvd0ZURVRNQkVHQTFVRQpBeE1LYldsdWFXdDFZbVZEUVRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTTYwCjd2SjNISTh3bG9XSWF2UkJneDRzWXROZGNLeFpiRlBWYk8rajRzemVVRlN1TllCUitXdmlHMzgyelQvKzVBNDgKT28wUHhib3g3djdRQXBua1RLTGNrT21xR05DN3RlOHNRVWtFemRSWGM3Y0x0Nmh4dHBrb2NraEtwK0ltM1BJQwpvQ2hzQmMwTWxWSmpiR25IRHJXZjFZa2JxaHQ1c3lBRTBEbVVGbXJ6MWdvdEgwN0d1Rkw2b0oyWjJCeXdmaERzCjh1cFB3djR6VzJLbXV4MjRwNUNnOEhVeTV4N2NmeFhQN2tsbkY5bVFRSDVqcXp6eDRBTnBWWURJTzA0QXZJK3cKQ1JXdGtlT3hpSmxXRkRFcjF0NURtcmp2YUlYcWdGN0ROYWlhMWtpcm94czFKUlFGM0xTK21va09iNmt2N1ZNVgpMa1dwZ0VMMG9DRnZybFMyT2k4Q0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUIwR0ExVWRKUVFXCk1CUUdDQ3NHQVFVRkJ3TUNCZ2dyQmdFRkJRY0RBVEFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFDRXhNYXVQUGczbjk1S1RsMEMrYVNoNzF5dHoxUHFXaExDdGJLU2s4T2EzRzhlMWU2OAoxVitNWVZDYjUxa25IT0YzRE03WkJBWllUYXB0MXRXelN6VnBCN2dkd2dNdU1Od3l0aDg2aXJ2eXFpbm96Q1FnCitBdHBsbDVkV0RsRVQzNFhpT05tTVVlYnVCd0JxZnNUeFArUE9IN2c0aEZpV1pTZEhRWTM4aFZFMWF5Wk92cjkKb2F1WXhpSVdEd3pzTkhCRXZaeTltVllNYnQ1dy9BYjZUWllQdTRiSVRIRE5mdlN6Y2h4aGVYNmt4VjRQQ3ByOQpSYmp0V1VDa0hqaWdJQ3AxclZhTDZGS0ZRUUxtUjFhZkxzRmxXZ1NTTzdsM1RHQUNINWZPSjR2cnBnOUtBTUk4ClNCUyt5TWsrRHRsSHljRG1QT3FBb2t3clBpUmpNVXF4N2NLeAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==`
	// token := `ZXlKaGJHY2lPaUpTVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUprWldaaGRXeDBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5elpXTnlaWFF1Ym1GdFpTSTZJbVJsWm1GMWJIUXRkRzlyWlc0dGFEazJOamNpTENKcmRXSmxjbTVsZEdWekxtbHZMM05sY25acFkyVmhZMk52ZFc1MEwzTmxjblpwWTJVdFlXTmpiM1Z1ZEM1dVlXMWxJam9pWkdWbVlYVnNkQ0lzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG5WcFpDSTZJbU5rTnpVMk5HRmhMV1kxTTJJdE1URmxOeTA1WkRoa0xUQTRNREF5TnpaaE1URTNaQ0lzSW5OMVlpSTZJbk41YzNSbGJUcHpaWEoyYVdObFlXTmpiM1Z1ZERwa1pXWmhkV3gwT21SbFptRjFiSFFpZlEuQ05OcDRRaTVQcFFLQVJuRkw1MGh5M3RWdm5vVVM2YWdwOVBwUGxIYXBmeThURGY4cEhBZm1IazIzTFRGUjlJU3k2TWVEY01VNmtJVFlNOVF3akFOYm9TNVNXekw5ZTRxNjM4NHFtdUN6Z3JVbnFfbEY2OWNyNUF5ejYtMHU0RnQ0aV8yWVVFU2xxb3ltT3o0NGg4cENtcXZ1MWRMQkdFazRMa2VhTzNucUVYWG00dHRoOE1zLTlpY294elk0WDdXZmxtMjFYV1JZbDh0anhzb2JLejVSbUdmeUNIaXNqcVF5WHpnQXRNdTRSbGVXWGR1U2tIRmJIRTZ3LU1aUGRqbWE4RmcyREFoTlJJbFc5U3BVckstT2dUZ1R6VDZCaUZSTmctU3pnU0x4U3pLWmRlX0FaRk41dl9nTS1QMHpQWFl3N0c5Q3d0X3hXemJlZ3FvbVRCYXVn`
	// kcli, err := NewContainerEngine(logrus.New(), host, cacert, token)
	// if err != nil {
	// 	fmt.Printf("\nerror occured when connecting kubernetes : %v", err)
	// }

	// fmt.Printf("\nKubeclient := %v", kcli)

	// c, err := kcli.CreateContainer(&CreateContainerOptions{Image: "nginx", ImageVersion: "1.13.6"})

	// if err != nil {
	// 	t.Log(err)
	// }

	// j, err := json.Marshal(c)

	// if err != nil {
	// 	t.Log(err)
	// }

	// fmt.Printf("%v", string(j))
}

// func testCreatePersistentVolume(t *testing.T) {
// 	kc := getKubernetesClient()
// 	pvo := &CreatePersistentVolumeOption{
// 		Path:         "/nfs/elasticshift",
// 		Server:       "10.10.3.128",
// 		MountOptions: []string{"hard", "nfsver=4.1"},
// 		provider:     NetworkFileShare,
// 		Name:         "pv-nfs-es",
// 		Capacity:     "5G",
// 	}
// 	res, err := kc.CreatePersistentVolume(pvo)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Create Persistent Volume %v", res)
// }

//func TestCreatePersistentVolumeHostDirectory(t *testing.T) {
//kc := getKubernetesClient()
//pvo := &CreatePersistentVolumeOption{
//Path:     "/Users/shahm/sandbox/es-nfs/hostpath",
//provider: HostLocalDirectory,
//}
//res, err := kc.CreatePersistentVolume(pvo)

//if err != nil {
//panic(err)
//}

//fmt.Printf("Create Persistent Volume %v", res)
//}

//func getKubernetesClient() KubernetesClient {

//	confBuf, err := base64.StdEncoding.DecodeString(confEncoded)
//	if err != nil {
//		fmt.Println("Error during kube initialization: ", err)
//	}

//	kco := &KubernetesClientOptions{KubeConfig: confBuf}
//	// kco := &KubernetesClientOptions{KubeConfigFile: "/Users/ghazni/.kube/config"}
//	//kco := &KubernetesClientOptions{KubeConfigFile: "$HOME/.kube/config"}
//	kc, err := NewKubernetesClient(kco)
//	if err != nil {
//		fmt.Printf("Error from test run:  %v ", err)
//	}

//	fmt.Printf("%v", kc)
//	r, _ := json.Marshal(kc)
//	fmt.Printf("The Marshal string %v", string(r))
//	return kc
//}
