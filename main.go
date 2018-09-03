package main

import (
	"encoding/json"
	"fmt"
	"github.com/appscode/go/log"
	"github.com/fatih/structs"
	"github.com/tamalsaha/go-oneliners"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type ObservableObject interface {
	GetResourceVersion() string
	GetGeneration() int64
	GetDeletionTimestamp() *metav1.Time
	GetLabels() map[string]string
	GetAnnotations() map[string]string
	//GetObservedGeneration() int64
	//GetObservedGenerationHash() string
}

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube/config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	kc := kubernetes.NewForConfigOrDie(config)
	deploys, err := kc.Apps().Deployments("kube-system").List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	var o1 interface{} = &deploys.Items[0]
	var o2 ObservableObject = o1.(ObservableObject)

	s2 := structs.New(o2)
	for _, f := range s2.Fields() {
		fmt.Println("**__ " + f.Name())
	}

	acc := o1.(metav1.Object)

	data, _ := json.MarshalIndent(o2, "", "  ")
	oneliners.FILE(string(data))
	oneliners.FILE(acc.GetNamespace())

	st := structs.New(o1)

	oneliners.FILE(st.Field("Status").Field("ObservedGeneration").Value().(int64))


	k, ok := st.FieldOk("xyz")
	k.IsZero()
	oneliners.FILE(ok)

	f2 := st.Field("Status").Field("ObservedGenerationHash")
	f2.IsZero()

	oneliners.FILE(f2.Value().(string))
}
