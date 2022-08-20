kubic
===

tui for k8s

## What is it.

`kubic` is a tui k8s explorer. 

The idea is to free from the need to keep namespaces, deployments, pods names in memory.

It's a good replacement for kubectl or more precisely it's an additional tool that takes over the most of work.

## Features

- doesn't require configuration for start, just kubernetes config
- vim mappings + arrows for navigation
- simple, sweet design powered by [Charm](https://charm.sh) libraries


## UI screenshots

<img src="assets/namespaces.png" width="700">
<img src="assets/deployments.png" width="700">
<img src="assets/pods.png" width="700">

Different color scheme example.

<img src="assets/theme.png" width="700">

## Requirements

- kubernetes `config` file

The best way is to keep it in `~/.kube/config` path. But you can set path on running.

## Installation

- Download archive from Releases
- Extract archive

```shell
chmod +x kubic
mv kubic ~/.local/bin
```

## Run

1. With default config.

 You can run `kubic` without any parameters if you have kubernetes config in `~/.kube/config` path nad need to use it.

```shell
kubic
```

2. With custom config.

You can run it with `--config` or shot `-c` parameter and set path to the config path. Or with environment variable `KUBIC_KUBERNETES_CONFIG_PATH`.

```shell
kubic -c /path/to/the/config/file
```

3. With custom color scheme.

If you want's to use different color scheme, put `style.json` in the same directory with `kubic` binary and run `kubic` without additional parameters. `kubic` will use this file automatically.

You can set path path to your json file with `--theme` parameter (short way `-t`) or with environment variable `KUBIC_THEME_FILE_PATH`.

```shell
kubic -t /path/to/the/json/style/file
```

Flags list:

| Short Flag | Long Flag | Environment Variable| Is Required | Type |
| ---   | --- | --- | --- | --- |
| -c | --config | KUBIC_KUBERNETES_CONFIG_PATH | False | string |
| -t | --theme | KUBIC_THEME_FILE_PATH | False | string |


## Customization

You can set your own color scheme with json file.

[Example](./assets/style.json)

Classes are predefined. Set your own color with hex.

```json
{
    "main-text": "#ffffff",
    "selected-text": "#61b0de",
    "inactive-text": "#616363",
    "tab-borders": "#109f93",
    "namespace-sign": "#eb24a9"
}
```

## Roadmap

- [x] k8s calls list
- [x] Namespaces list (kubectl get namespace)
- [x] Deployments list (kubectl get deployments --namespace=ns)
- [x] Pods list (kubectl get pods --namespace=ns)
- [ ] Deployment info (kubectl get deployment deployment-name --namespace=ns -o=json)
- [ ] Pod info (kubectl get pods pod-name --namespace=ns -o=json)
- [ ] Pod logs (kkubectl logs pod-name --namespace=ns)

***

Powered by [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="200"></a>