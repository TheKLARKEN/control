until $([ $(sudo kubectl get nodes|grep Ready|wc -l) -eq {{ .MachineCount }} ]); do printf '.'; sleep 5; done
