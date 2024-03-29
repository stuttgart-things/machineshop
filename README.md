# stuttgart-things/machineShop

git based CLI interface for managing configuration as code

## FEATURES
* RENDER TEMPLATES w/ DEFAULTS AND PARAMETERS (RENDER)
* INSTALL MULTIPLE BINARIES FROM WEB SOURCES AT ONCE/IN PARALLEL (INSTALL)
* RENDER + EXECUTE MULTIPLE SCRIPTS (INSTALL)
* RETRIEVE SECRETS FROM VAULT (GET)

## STATUS
* 60% FEATURES DONE

## INSTALLATION

<details><summary><b>BY RELEASE</b></summary>

```bash
# LINUX x86_64
VERSION=0.1.48
wget https://github.com/stuttgart-things/machineShop/releases/download/${VERSION}/machineShop_Linux_x86_64.tar.gz
tar xvfz machineShop_Linux_x86_64.tar.gz
sudo mv machineShop /usr/bin/machineShop
rm -rf LICENSE README.md
sudo chmod +x /usr/bin/machineShop
machineShop version
```

</details>


## USAGE EXAMPLES

### PUSH

<details><summary><b>GIT</b></summary>

```bash
export MINIO_ACCESS_KEY=sthings
export MINIO_SECRET_KEY=<PASSWORD>
export MINIO_ADDR=artifacts.automation.sthings-vsphere.labul.sva.de
export MINIO_SECURE=true

machineShop push \
--target minio \
--source pod.yaml \
--destination manifests:pod-example.yaml # <BUCKET>:<OBECTNAME>
```

</details>


### RENDER

<details><summary><b>GIT</b></summary>

```bash
machineShop render --source git \
--git https://github.com/stuttgart-things/stuttgart-things.git \
--defaults packer/environments/labul-vsphere.yaml \
--template packer/os/ubuntu23-vsphere.pkr.tpl.hcl \
--output stdout
```

</details>

<details><summary><b>LOCAL</b></summary>

```bash
machineShop render \
--source local \
--template ../golang/machineShop/tests/template-square.yaml \
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
machineShop get --path apps/data/scr:password | tail -n +8

machineShop get --path apps/data/scr:password --output file --destination /tmp/password.txt

machineShop get --path kubeconfigs/data/dev21:kubeconfig --output file --destination /tmp/dev211 --b64 true
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
