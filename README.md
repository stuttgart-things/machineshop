# stuttgart-things/machineshop

git based CLI interface for managing configuration as code

## TASKS

```bash
task: Available tasks for this project:
* branch:              Create branch from main
* build:               Build code
* build-image:         Build container image
* commit:              Commit + push code into branch
* delete-branch:       Delete branch from origin
* install:             Install
* ko:                  Build image w/ KO
* lint:                Lint
* pr:                  Create pull request into main
* release:             Relase binaries
* run:                 Run
* tag:                 Commit, push & tag the module
* test:                Test
* tests:               Built cli tests
```

TASK EXAMPLES

```bash
task run # will output build version
task run CMD=get PARAMETERS=--system=sops # will run with build command get + parameters
task release TAG=2.2.9 # will release bins with version 2.2.9
```


## FEATURES
* RENDER TEMPLATES w/ DEFAULTS AND PARAMETERS (RENDER)
* INSTALL MULTIPLE BINARIES FROM WEB SOURCES AT ONCE/IN PARALLEL (INSTALL)
* RENDER + EXECUTE MULTIPLE SCRIPTS (INSTALL)
* RETRIEVE SECRETS FROM VAULT (GET)

## DEPLOYMENT

<details><summary><b>BINARY BY RELEASE</b></summary>

```bash
# LINUX x86_64
VERSION=v1.9.0
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

<details><summary><b>CREATE BRANCH</b></summary>

```bash
task branch
```

</details>

<details><summary><b>CREATE PULL-REQUEST/MERGE</b></summary>

```bash
task pr
```

</details>

<details><summary><b>BUILD RELEASE</b></summary>

```bash
task release TAG=v1.8.0 # EXAMPLE VERSION
```

</details>

<details><summary><b>BUILD CONTAINER-IMAGE w/ KO</b></summary>

```bash
task ko TAG=v1.9.0 # EXAMPLE VERSION
```

</details>

## USAGE EXAMPLES

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
--group stuttgart-things
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

### PUSH

push things to targets

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

<details><summary><b>VAULT-REQUIREMENT: VAULT APPROLE EXPORTS</b></summary>

```bash
export VAULT_NAMESPACE=root
export VAULT_ROLE_ID=1d42d7e7-8c14-e5f9-801d-b3ecef416616
export VAULT_SECRET_ID=623c991f-dd76-c437-2723-bb2ef5b02d87
export VAULT_ADDR=https://≤VAULT_ADDR>[:8200]
```

</details>

<details><summary><b>SOPS-REQUIREMENT: AGE_KEY EXPORTS</b></summary>

```bash
export SOPS_AGE_KEY=AGE-...
# or
export SOPS_AGE_KEY_FILE=home/sthings/projects/golang/sops/sops.key
```

</details>

<details><summary><b>GET SECRET VALUE BY PATH</b></summary>

```bash
machineshop get --path apps/data/scr:password | tail -n +8

machineshop get --path apps/data/scr:password --output file --destination /tmp/password.txt

machineshop get --path kubeconfigs/data/dev21:kubeconfig --output file --destination /tmp/dev211 --b64 true
```

</details>

<details><summary><b>GET SECRET VALUE BY PATH</b></summary>

```bash
machineshop get --system=sops --path=/home/sthings/projects/golang/sops/bla.yaml:age | tail -n +11
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
