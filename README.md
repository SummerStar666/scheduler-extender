# scheduler-extender
kubernetes自定义scheduler-extender。

## 调度器工作原理

1. 默认调度器根据指定的参数启动（我们使⽤ kubeadm 搭建的集群，启动配置⽂件位于 /etc/kubernetes/manifests/kube-schdueler.yaml） 
2. watch apiserver，将 spec.nodeName 为空的 Pod 放⼊调度器内部的调度队列中 
3. 从调度队列中 Pop 出⼀个 Pod，开始⼀个标准的调度周期 
4. 从 Pod 属性中检索“硬性要求”（⽐如 CPU/内存请求值，nodeSelector/nodeAffinity），然后过滤阶段发⽣，在该阶段计算出满⾜要求的节点候选列表 
5. 从 Pod 属性中检索“软需求”，并应⽤⼀些默认的“软策略”（⽐如 Pod 倾向于在节点上更加聚拢或分散），最后，它为每个候选节点给出⼀个分数，并挑选出得分最⾼的最终获胜者 
6. 和 apiserver 通信（发送绑定调⽤），然后设置 Pod 的 spec.nodeName 属性以表示将该 Pod 调度到的节点。

## 自定义extender工作逻辑

### filter
filter函数调用podFitsOnNode函数，来判断当前pod是否适合放在这个节点上。判断依据为，一个随机数对5求余数，这个余数是否等于3。如果等于3，则证明这个pod可以放在该节点上。
```go
func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	var failReasons []string
	judge := (rand.Intn(100)%5 == 3)
  
	if judge {
		log.Printf("pod %v/%v fits on node %v\n", pod.Name, pod.Namespace, node.Name)
		return true, nil, nil
	}
	log.Printf("pod %v/%v does not fit on node %v\n", pod.Name, pod.Namespace, node.Name)

	failures := "It's not fits on this node."
	failReasons = append(failReasons, failures)
  
	return false, failReasons, nil
}
```

### prioritize
在prioritize函数中，每个节点按照在队列中的顺序打分，分数依次升高。
```go
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
```
