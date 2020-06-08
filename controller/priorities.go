package controller

import (
	"log"

	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

const (
	PrioMsg = "pod %v/%v gets score %v\n"
)

//打分，排序
func prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	pod := args.Pod
	nodes := args.Nodes.Items

	hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes))
	for i, node := range nodes {
		score := i+1
		log.Printf(PrioMsg, pod.Name, pod.Namespace, score)
		hostPriorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}

	return &hostPriorityList
}
