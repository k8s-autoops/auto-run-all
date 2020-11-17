# auto-run-all

自动在所有 Pod 内跑一段脚本，用于执行某些维护任务，比如删除一些棘手的日志文件

默认使用 `sh`，按需要把要执行的脚本挂载到 `/autoops-data/auto-run-all/script.sh` 这个位置

## Usage

Create namespace `autoops` and apply yaml resources as described below.

```yaml
# create serviceaccount
apiVersion: v1
kind: ServiceAccount
metadata:
  name: auto-run-all
  namespace: autoops
---
# create clusterrole
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: auto-run-all
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["list"]
  - apiGroups: [""]
    resources: ["pods/exec"]
    verbs: ["create"]
---
# create clusterrolebinding
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: auto-run-all
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: auto-run-all
subjects:
  - kind: ServiceAccount
    name: auto-run-all
    namespace: autoops
---
# create cronjob
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: auto-run-all
  namespace: autoops
spec:
  schedule: "*/5 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccount: auto-run-all
          containers:
            - name: auto-run-all
              image: autoops/auto-run-all
          restartPolicy: OnFailure
```

## Credits

Guo Y.K., MIT License
