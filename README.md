# stuttgart-things/machineshop

git based CLI interface for managing configuration as code

## FEATURES
* RENDER TEMPLATES w/ DEFAULTS AND PARAMETERS (RENDER)
* INSTALL MULTIPLE BINARIES FROM WEB SOURCES AT ONCE/IN PARALLEL (INSTALL)
* RENDER + EXECUTE MULTIPLE SCRIPTS (INSTALL)
* RETRIEVE SECRETS FROM VAULT (GET)

## INSTALLATION

<details><summary><b>BY RELEASE</b></summary>

```bash
# LINUX x86_64
VERSION=v1.7.0
wget https://github.com/stuttgart-things/machineshop/releases/download/${VERSION}/machineshop_Linux_x86_64.tar.gz
tar xvfz machineshop_Linux_x86_64.tar.gz
sudo mv machineshop /usr/bin/machineshop
rm -rf LICENSE README.md
sudo chmod +x /usr/bin/machineshop
machineshop version
```

</details>

## DEV

<details><summary><b>BUILD RELEASE</b></summary>

```bash
task release TAG=v1.8.0 # EXAMPLE VERSION
```

</details>

<details><summary><b>BUILD CONTAINER-IMAGE w/ KO</b></summary>

```bash
task ko
```

</details>

## USAGE EXAMPLES

### CREATE

creates things on github

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

### PUSH

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


### RENDER

<details><summary><b>GIT</b></summary>

```bash
machineshop render --source git \
--git https://github.com/stuttgart-things/stuttgart-things.git \
--defaults packer/environments/labul-vsphere.yaml \
--template packer/os/ubuntu23-vsphere.pkr.tpl.hcl \
--output stdout
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


<details><summary><b>GET</b></summary>

### REQUIREMENT: VAULT APPROLE EXPORTS

```bash
export VAULT_NAMESPACE=root
export VAULT_ROLE_ID=1d42d7e7-8c14-e5f9-801d-b3ecef416616
export VAULT_SECRET_ID=623c991f-dd76-c437-2723-bb2ef5b02d87
export VAULT_ADDR=https://≤VAULT_ADDR>[:8200]
```

### GET SECRET VALUE BY PATH
```
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
