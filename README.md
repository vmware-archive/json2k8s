# json2k8s

json2k8s extracts Kubernetes resources [1] from JSON files.

```bash
$ cat examples/file1.json | json2k8s
$ json2k8s examples/file1.json examples/file2.json
$ jsonnet examples/file1.json | json2k8s | kapp deploy -a app1 -f- --yes
```

Features:

- extracts Kubernetes resources from arbitrarily nested JSON structures
- prints result as single v1.List resource
- no de-duplication of resources is done (expected to be done by downstream tools e.g. kapp)

[1] Kubernetes resource is considered to be a map with `kind` and `apiVersion` keys present

## Motivation

Above functionality can be useful for defining applications in Jsonnet in a more composable manner. [kubecfg](https://github.com/bitnami/kubecfg/) bakes this in; however, since k14s tools are written to be composable this was split off into its own binary. It can run on [jsonnet](https://jsonnet.org/) result and ultimately passed to [k14s/kapp](https://get-kapp.io).
