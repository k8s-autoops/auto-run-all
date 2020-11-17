package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
)

const (
	ScriptPath = "/autoops-data/auto-run-all/script.sh"
)

func exit(err *error) {
	if *err != nil {
		log.Println("exited with error:", (*err).Error())
		os.Exit(1)
	} else {
		log.Println("exited")
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lmsgprefix)

	var err error
	defer exit(&err)

	var buf []byte
	if buf, err = ioutil.ReadFile(ScriptPath); err != nil {
		return
	}

	var cfg *rest.Config
	if cfg, err = rest.InClusterConfig(); err != nil {
		return
	}

	var client *kubernetes.Clientset
	if client, err = kubernetes.NewForConfig(cfg); err != nil {
		return
	}

	var pods *corev1.PodList
	if pods, err = client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{}); err != nil {
		return
	}

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			execute(cfg, client, pod, container, bytes.NewReader(buf))
		}
	}
}

func execute(cfg *rest.Config, client *kubernetes.Clientset, pod corev1.Pod, container corev1.Container, stdin io.Reader) {
	var err error
	log.Printf("Pod: %s/%s (Container: %s)", pod.Namespace, pod.Name, container.Name)
	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		Container: container.Name,
		Command:   []string{"sh"},
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
	}, scheme.ParameterCodec)
	var exec remotecommand.Executor
	if exec, err = remotecommand.NewSPDYExecutor(cfg, "POST", req.URL()); err != nil {
		log.Printf("Failed to create executor: %s", err.Error())
		return
	}
	if err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}); err != nil {
		log.Printf("Failed to execute script: %s", err.Error())
		return
	}
}
