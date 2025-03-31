# stuttgart-things/machineshop

<div align="center">
  <p>
    <img src="https://github.com/stuttgart-things/docs/blob/main/hugo/sthings-train.png" alt="sthings" width="350" />
  </p>
  <p>
    <strong>[məˈʃiːnʃɒp]</strong>- git based CLI interface for managing configuration as code

  </p>
</div>

## FEATURES
* RENDER TEMPLATES w/ DEFAULTS AND PARAMETERS (RENDER)
* INSTALL MULTIPLE BINARIES FROM WEB SOURCES AT ONCE/IN PARALLEL (INSTALL)
* RENDER + EXECUTE MULTIPLE SCRIPTS (INSTALL)
* RETRIEVE SECRETS FROM VAULT (GET)

## DEPLOYMENT

<details><summary><b>FROM RELEASE-ARCHIVE</b></summary>

```bash
# LINUX x86_64
VERSION=v2.6.10
wget https://github.com/stuttgart-things/machineshop/releases/download/${VERSION}/machineshop_Linux_x86_64.tar.gz
tar xvfz machineshop_Linux_x86_64.tar.gz
sudo mv machineshop /usr/bin/machineshop
rm -rf LICENSE README.md
sudo chmod +x /usr/bin/machineshop
machineshop version
```

</details>

<details><summary><b>CONTAINER IMAGE</b></summary>

```bash
# RUN COMMAND
sudo nerdctl run ghcr.io/stuttgart-things/machineshop/machineshop-9c3178088556daa12a17db5edcc6b5b7:1.9.10 version
```

```bash
# JUMP INTO SHELL
nerdctl run -it --entrypoint bash \
ghcr.io/stuttgart-things/machineshop/machineshop-9c3178088556daa12a17db5edcc6b5b7:1.9.10
```

</details>

## DEV

<details><summary>ALL TASKS</summary>

```bash
task: Available tasks for this project:
* branch:              Create branch from main
* build:               Build code
* build-image:         Build container image
* build-ko:            Build image w/ KO
* check:               Run pre-commit hooks
* commit:              Commit + push code into branch
* dagger-ko:           Build image w/ ko
* delete-branch:       Delete branch from origin
* goreleaser:          Release bins w/ goreleaser
* install:             Install
* lint:                Lint
* pr:                  Create pull request into main
* release:             Release
* release-image:       Release image
* run:                 Run
* switch-local:        Switch to local branch
* switch-remote:       Switch to remote branch
* tag:                 Commit, push & tag the module
* tasks:               Select a task to run
* test:                Test
* test-install:        Test crossplame modules
* test-run:            Build Bin & Test Run Command
* test-version:        Test version cmd
* tests:               Built cli tests
```

</details>

<details><summary>SELECT TASK</summary>

```bash
task=$(yq e '.tasks | keys' Taskfile.yaml | sed 's/^- //' | gum choose) && task ${task}
```

</details>

## USAGE

### CREATE

creates things on github

<details><summary><b>REPOSITORY</b></summary>

```bash
export GITHUB_TOKEN=<GITHUB_TOKEN>

machineshop create \
--kind repo \
--group stuttgart-things \
--repository machineshop2 \
--message "test repository - machineshop" \
--private true
```

</details>

<details><summary><b>BRANCH</b></summary>

```bash
export GITHUB_TOKEN=<GITHUB_TOKEN>

machineshop create \
--kind branch \
--branch hello \
--repository machineshop \
--group stuttgart-things \
--files "Dockerfile:Dockerfile" \
```

</details>

<details><summary><b>PULL-REQUEST</b></summary>

```bash
export GITHUB_TOKEN=<GITHUB_TOKEN>

machineshop create \
--kind pr \
--title test2 \
--branch hello \
--repository machineshop \
--group stuttgart-things \
--labels "release,deploy" # optional
```

</details>

<details><summary><b>LABELS (ON PULL-REQUEST)</b></summary>

all existing labels will be overwritten by the specified / to be upodated.

```bash
export GITHUB_TOKEN=<GITHUB_TOKEN>

machineshop create \
--kind labels \
--group stuttgart-things \
--repository kaeffken \
--id 58 \
--labels app1,deploy
```

</details>

<details><summary><b>MERGE</b></summary>

```bash
export GITHUB_TOKEN=<GITHUB_TOKEN>

machineshop create \
--kind merge \
--group stuttgart-things \
--repository stuttgart-things \
--message "test" \
--merge rebase \
--id 243
```

</details>

### RUN

<details><summary><b>RUN</b></summary>

```bash
# PROFILE FROM LOCAL
sudo machineshop run \
--scripts golang \
--source local \
--profile run.yaml
```

```bash
# DEFAULT PROFILE BY URL
sudo machineshop run \
--scripts buildx,onefetch
```

</details>

### INSTALL

<details><summary><b>INSTALL</b></summary>

```bash
sudo machineshop install \
--profile machineShop/binaries.yaml \
--binaries "sops,flux"
```

</details>

### PUSH

push things to targets

<details><summary><b>SET IPS TO CLUSTER (w/ CLUSTERBOOK)</b></summary>

```bash
machineshop push \
--target=ips \
--destination=clusterbook.rke2.sthings-vsphere.labul.sva.de:443 \
--artifacts="10.31.103.9;10.31.103.10" \
--assignee=app1
```

</details>

<details><summary><b>MINIO (BUCKET)</b></summary>

```bash
export MINIO_ACCESS_KEY=sthings
export MINIO_SECRET_KEY=<PASSWORD>
export MINIO_ADDR=artifacts.automation.sthings-vsphere.labul.sva.de
export MINIO_SECURE=true

machineshop push \
--target minio \
--source pod.yaml \
--destination manifests:pod-example.yaml # <BUCKET>:<OBECTNAME>
```

</details>

<details><summary><b>MS TEAMS</b></summary>

```bash
WEBHOOK_URL=https://365sva.webhook...

machineshop push \
--target teams \
--source "hello from machineshop cli" \
--destination ${WEBHOOK_URL} \
--color blue
```

</details>

<details><summary><b>HOMERUN-DEMO</b></summary>

```bash
HOMERUN_URL=https://homerun.homerun-dev.sthings-vsphere.labul.sva.de/generic

machineshop push \
--destination ${HOMERUN_URL} \
--target homerun-demo \
--source profiles/homerun.yaml
```

</details>


<details><summary><b>HOMERUN</b></summary>

```bash
HOMERUN_URL=https://homerun.homerun-dev.sthings-vsphere.labul.sva.de/generic

machineshop push \
--destination ${HOMERUN_URL} \
--target homerun \
--title "hello" \
--system shell \
--message "test sdfsdfslkljh" \
--tags "shell;linux" \
--author "machineshop" \
--severity "INFO"

machineshop push \
--destination ${HOMERUN_URL} \
--target homerun \
--title "hello" \
--system shell \
--message "test sdfsdfslkljh" \
--tags "shell;linux" \
--author "machineshop" \
--severity "INFO" \
--assignee "patrick.hermann" \
--assigneeUrl "patrick.hermann@sva.de"  \
--artifacts "INFO" \
--url "https://github.com/stuttgart-things/stuttgart-things/actions/runs/10639438939"
```

</details>


### RENDER

render things from templates from various input sources

<details><summary><b>EXAMPLE TEMPLATE</b></summary>

```yaml
---
runs:
  packagePublishHelmChart:
    # FLAT VALUE
    name: package-publish-{{ .chartName }}

# LOOP OVER LIST
{{ range .food }}
- {{ . }}{{ end }}

# RANDOM ELEMENT FROM EXISTING LIST
favoriteFood: {{ .RANDOMfood }}
cpu: {{ .vmConfig_l_cpu }}
ram: {{ .vmConfig_m_ram }}
```

</details>

<details><summary><b>EXAMPLE MULTIKEY-TEMPLATE</b></summary>

```yaml
---
template:
  nfsCsi: |
    apiVersion: kustomize.toolkit.fluxcd.io/v1
    kind: Kustomization
    metadata:
      name: nfs-csi
      namespace: {{ .namespace }}
    spec:
      interval: {{ .interval }}
# ...
longhorn: |
    apiVersion: kustomize.toolkit.fluxcd.io/v1
    kind: Kustomization
    metadata:
      name: longhorn
      namespace: {{ .namespace}}
#...
```

</details>

<details><summary><b>EXAMPLE SOPS DECRYPTED (VALUES) FILE</b></summary>

```yaml
vm_ssh_user: ENC[AES256_GCM,data:2nlFCn5/qA==,iv:+AvMEg1RHFlBqRRNloXNTxTEaUvq1x1tNM4S2liE9is=,tag:2+DvYvHtNSSojg9N7yTKbQ==,type:str]
vm_ssh_password: ENC[AES256_GCM,data:XGQ+GjNqnhA=,iv:UIO5+4vOiGWOlonBKF4wb0n2Pj9VBngieMfcBDmQdXM=,tag:O44ZDEE7nBtECPRxcQsSuw==,type:str]
sops:
    kms: []
    gcp_kms: []
    azure_kv: []
    hc_vault: []
    age:
        - recipient: age1g438n4lx6h7x7u42q652e9ygzrkkwlul49e8zsmsrfmxm9k3tvcsykhff4
          enc: |
            -----BEGIN AGE ENCRYPTED FILE-----
            YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSAraTYrTExhMGV5Y1lYb0g5
            OWswNUprbDdobTk0N0UyQUZiZmxRWS9wdkRjCmRZMWw3dE1VNjZ5M0xRaGp3NmQ1
            UVFRQ0hhY3pRV0dmY2YyQ0lFRjFqSk0KLS0tIGdrY1FXYnJNYy8wVW5XcHZENkhV
            VUxGV09pVllCYXU0dXhlZFdDWXBMRmsKs55x9DeiqsRjLSm+U+BVdsJ6dLeNqeSE
            xJ/3GQpy/MmyARUzayTSOuzu8URemMaAh7FQbxTf8V7AnMM6Lv+sHQ==
            -----END AGE ENCRYPTED FILE-----
    lastmodified: "2025-01-01T09:47:30Z"
    mac: ENC[AES256_GCM,data:Nlqp8wiKQbGTzP3UuPlNMp7rmyEcyiGQVsFushGtKpWOAiSudcvZEKPkdFGQz3KJQXVg7KCBK1RqYTo0ZIw6Fwr0cbIEwHFekVSnpEu5p2H0JVDtlpO+clkEZqNi/HVGnIF16cYpNGrnI74tD4DaaS4HAF04OcvAHHBAtmTQhuk=,iv:4q4KQWomb3PxQa8bTHrdNx+02aLeQ/P+RIpf8rifRrc=,tag:UgWmC56SbWe2KiG4IsaJ6A==,type:str]
    pgp: []
    unencrypted_suffix: _unencrypted
    version: 3.8.1
```

</details>

<details><summary><b>EXAMPLE DEFAULTS FILE</b></summary>

```yaml
---
chartName: helloHelm
food:
  - schnitzel
  - apple
  - hamburger

vmConfig:
  m:
    cpu: 6
    ram: 8192
  l:
    cpu: 8
    ram: 10240
```

</details>

<details><summary><b>LOCAL</b></summary>

```bash
machineshop render \
--source local \
--template ../golang/machineshop/tests/template-square.yaml \
--brackets square \
--output stdout \
--defaults /home/sthings/projects/stuttgart-things/packer/environments/labul-pve.yaml
```

</details>

<details><summary><b>GIT</b></summary>

```bash
machineshop render --source git \
--git https://github.com/stuttgart-things/stuttgart-things.git \
--defaults packer/environments/labul-vsphere.yaml \
--template packer/os/ubuntu23-vsphere.pkr.tpl.hcl \
--output stdout
```

</details>

<details><summary><b>RENDER FROM MULTIKEY YAML TEMPLATE</b></summary>

```bash
machineshop render \
--source local \
--template tests/infra.yaml \
--output stdout \
--kind multikey \
--key longhorn \
--defaults tests/default.yaml
```

</details>

<details><summary><b>RENDER w/ SECRETS FROM SOPS</b></summary>

```bash
export SOPS_AGE_KEY=AGE-SECRET-KEY-1T22K05UTR...

machineshop render \
--source local \
--template tests/infra.yaml \
--output stdout \
--secrets secrets1.yaml
```

</details>


### DELETE

delete things on git(hub)

<details><summary><b>BRANCH</b></summary>

```bash
export GITHUB_TOKEN=<GITHUB_TOKEN>

machineshop delete \
--kind branch \
--branch hello \
--repository stuttgart-things \
--group stuttgart-things
```

</details>

<details><summary><b>FILES</b></summary>

```bash
export GITHUB_TOKEN=<GITHUB_TOKEN>

machineshop delete \
--kind files \
--branch main \
--repository stuttgart-things \
--group stuttgart-things \
--files ".github/workflows/lint-k8s-manifests.yaml" \
--user patrick-hermann-sva
```

</details>

### GET

get things from systems

<details><summary><b>GET IPS FROM CLUSTERBOOK</b></summary>

```bash
machineshop get \
--system=ips \
--destination=clusterbook.rke2.sthings-vsphere.labul.sva.de:443 \
--path=10.31.103 \
--output=2
```

</details>

<details><summary><b>VAULT-REQUIREMENT: VAULT APPROLE EXPORTS</b></summary>

```bash
export VAULT_NAMESPACE=root
export VAULT_ROLE_ID=1d42d7e7-8c14-e5f9-801d-b3ecef416616
export VAULT_SECRET_ID=<SECRET>
export VAULT_ADDR=https://≤VAULT_ADDR>[:8200]
```

</details>

<details><summary><b>SOPS-REQUIREMENT: AGE_KEY EXPORTS</b></summary>

```bash
export SOPS_AGE_KEY=AGE-...
```

</details>

<details><summary><b>GET SOPS SECRET VALUE BY PATH</b></summary>

```bash
export SOPS_AGE_KEY=AGE-SECRET-KEY-1T22K0..
machineshop get --system=sops --path=/home/sthings/projects/golang/sops/bla.yaml:password | tail -n +1
```

</details>

<details><summary><b>GET VAULT SECRET VALUE BY PATH</b></summary>

```bash
machineshop get --path apps/data/scr:password | tail -n +8

machineshop get --path apps/data/scr:password --output file --destination /tmp/password.txt

machineshop get --path kubeconfigs/data/dev21:kubeconfig --output file --destination /tmp/dev211 --b64 true
```

</details>



## LICENSE

<details><summary><b>APACHE 2.0</b></summary>

Copyright 2023 patrick hermann.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

</details>

Author Information
------------------
Patrick Hermann, stuttgart-things 05/2023
