# Edgex Helm Chart Encryption
## _Introduction_

[![N|Solid](data:image/png;base64,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)](https://nodesource.com/products/nsolid)

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

The aim of this repo is to provide a basic EdgeX deployable helm chart, with a secret encryption step for the values.yaml file

## Prerequisites

* Install k3s:

```
curl -sfL https://get.k3s.io | sh -
```

* Install Homebrew:

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

If this step cause error (i.e. curl 56), just check the internet connection, remove the homebrew folder from _/home_ and reinstall it again

* Install k9s:

```
brew install k9s
```

* Install Helm:

```
brew install helm
```

## Helm Chart Strucutre

📦tim-edgex
 ┣ 📂templates
 ┃ ┣ 📂core-command
 ┃ ┃ ┣ 📜edgex-core-command-deployment.yaml
 ┃ ┃ ┗ 📜edgex-core-command-service.yaml
 ┃ ┣ 📂core-data
 ┃ ┃ ┣ 📜edgex-core-data-deployment.yaml
 ┃ ┃ ┗ 📜edgex-core-data-service.yaml
 ┃ ┣ 📂core-metadata
 ┃ ┃ ┣ 📜edgex-core-metadata-deployment.yaml
 ┃ ┃ ┗ 📜edgex-core-metadata-service.yaml
 ┃ ┣ 📂redis
 ┃ ┃ ┣ 📜edgex-redis-deployment.yaml
 ┃ ┃ ┗ 📜edgex-redis-service.yaml
 ┃ ┣ 📂support-scheduler
 ┃ ┃ ┣ 📜edgex-support-scheduler-deployment.yaml
 ┃ ┃ ┗ 📜edgex-support-scheduler-service.yaml
 ┃ ┗ 📜configmap.yaml
 ┣ 📜Chart.yaml
 ┣ 📜README.md
 ┗ 📜<secrets.>values.yaml

Where we have:

* chart.yaml: Inform Helm that the content of the chart represents a Chart
* <secrets.>values.yaml: Contains environmental and custom variables that should be not available to external users
* templates: Contains an optional config_map.yaml file, used to map variables and a list of subfolders to be deployed, each one representing a single POD. 
* Each POD folder will contains:
    * \*.deployment.yaml: Contains information about how to deploy the service
    * \*.service.yaml: Contains specifications about the container
    

## Deployment

First of all check that the kubernetes service is up and running correctly:

* To check:

```
sudo service k3s status
```

* To start or restart:

```
sudo service k3s start
--or--
sudo service k3s restart
```

* To kill or stop:

```
cd /usr/local/bin && ./k3s-killall.sh
```

Then give the right permission to the kubeconfig file:

```
sudo chmod +r /etc/rancher/k3s/k3s.yaml
```

Finally, from a shell, type the following commands for:

* Deploying: 

```
helm install <distro-name> PathToHelmChart --kubeconfig /etc/rancher/k3s/k3s.yaml
```

* Upgrading:

```
helm upgrade <distro-name> PathToHelmChart --kubeconfig /etc/rancher/k3s/k3s.yaml
```

* Deleting:

```
helm delete <distro-name> --kubeconfig /etc/rancher/k3s/k3s.yaml
```

## Check PODs status:

At this point it is possible to chek the architecture status by typing the following command:

```
k9s --kubeconfig /etc/rancher/k3s/k3s.yaml
```

The following output should appear, if everything is working correctly:

[![N|Solid](https://miro.medium.com/max/1838/1*8OdadtRuO4Q9hbYEeKQjwg.png)](https://nodesource.com/products/nsolid)

## Encryption

The Helm security plug-in helps to define the _values.yaml_ file, encrypting sensible data. These encrypted data can be distributed and stored in the code version management tool at will, without worrying about the exposure of sensitive information. Here is an example of an encrypted values file:

```
#ENC[AES256_GCM,data:IHAqGPYHlUdD2+xSn5ZcYCo=,iv:1KKx8l1zl41LuNYcKw3biXm0vx+vjAeA7wdnNHYjQ6Y=,tag:MmWG4SIeXPt0o0HOHGtJeQ==,type:comment]
registry:
    url: ENC[AES256_GCM,data:sYON9+wBDq9jcmhy8iUaITpIjApLbys=,iv:z/ITKkJp2rS/jMyvxghweA+7W0QlZ98PR+4gDGhX+WI=,tag:TLbCkMcIfW30/80FpaozoA==,type:str]
    username: ENC[AES256_GCM,data:5Ju2bxk=,iv:hxRUoi0lViW7chOQTiyZyt4nGMS5V5YZyFNf19LmvpA=,tag:lb830A5pnZ4bI0HUosyc7Q==,type:str]
    password: ENC[AES256_GCM,data:1WmPZCSlzGbSn2LqMc7DmHjQKTjsaPUdn9nMKvbl+KIJ451EPkV0s3dqF+NVZ5E+T6reYN/lY7Ok3VmGvboVAbFs1IYzn7KenbGLMGgCT+JhUFaYz16TeGGsyWDk6YcIIw/XzR6lTjilpHF+DuZuepOyiAnCO0Q5k4aux2lICQh6P8mOezt8flP9/blnFGVZhaaE5r5vT6hsaQbsy7Rnk2lP926xT8NWcaXR85AleRvevQ/zwFIFjjk=,iv:ivA6U20LCHOoR9WGSmuvlJdhnYx/ZC8Pw9czMjNrrlI=,tag:wX3SiiBv8OjSZkpbCznDZw==,type:str]

jenkins:
    master:
        JCasC:
            configScripts:
                credentials-config: ENC[AES256_GCM,data:3dN7KBW...Ov/rsUA=,tag:qlfV/0x0vr45JxMYM1UdMQ==,type:str]
sops:
    kms: []
    gcp_kms: []
    lastmodified: '2019-07-10T06:21:36Z'
    mac: ENC[AES256_GCM,data:sKsL25V5yci+oD1PpfA5fU6zE7YCc6Sxg7myE4eqoDcA+guG8gUg4Hcj5yAB4APBq3+KtPIXoF0hNHVYZOOYqZXQrMpO0jASjWHmLAFTUb6FE6xOtb4mP3FBk8W6Km7TfNz3Te8WW4nsb/+c0WmFSQnIolaeXgbbZhZ23x+V9g=1,iv:Oha7rwD2y3xCc+UnI+xXwrnFByMhNJkF84TiYq4/LsWI=,tag:W3e9ox2G9QL5jQEV0VwGA==,type:str]
    pgp:
    -   created_at: '2019-07-10T06:21:34Z'
        enc: |
            -----BEGIN PGP MESSAGE-----

            hQEMA9Q2nDmrg55qAQf/aXiC7EXZlP5OZDrH3clCb0I9uqP8eNhVgAzqyfSaajGB
            ...
            =h7fE
            -----END PGP MESSAGE-----
        fp: AD331C18082B4669992805DDCB8EA0C7BC44A464
    unencrypted_suffix: _unencrypted
    version: 3.0.3
```

It can be seen that the encryption operation is only carried out for the values corresponding to all keys in the yaml file, and the true key is retained. Compared with the way of overall encryption for the whole yaml file, the encrypted file still retains strong readability, which makes it possible to maintain the encrypted yaml file.

When deploying the application, it is sufficient to directly pass the encrypted values file to the helm command, and the helm secrets plug-in will automatically decrypt the encrypted values without any human intervention.

### Install the Helm secrets plug-in

Before using helm secrets plug-in, first ensure that the plug-in is installed in the local helm:

```
helm plugin install https://github.com/futuresimple/helm-secrets
```

After installation, it is possible to check by typing:

```
helm plugin list

--OUTPUT--

NAME       VERSION    DESCRIPTION
secrets    2.0.2      This plugin provides secrets values encryption for Helm charts secure storing
```

### Introducing SOPs commands

Since the Helm secrets plug-in, does not have any encryption/decryption capabilities, it is necessary to pass through an external tool that provides these capabilities.

SOPs is an open source text editing tool developed by Mozilla. It supports editing files in yaml, JSON, env, ini and binary text formats, and encrypts/decrypts the edited files by using encryption methods such as AWS, kms, GCP kms, azure key vaul or PGP.

For now it is sufficient to understand that SOPs command is an essential dependency to ensure that the Helm secrets pulg-in can work normally. These commands are automatically installed by the Helm secrets plug-in, so after plug-in installation it is possible to perform a check by asking for SOPs version information:

```
sops -v
```

Of course, the tool can also be installed manually by typing:

```
brew install sops
```

### Introduction to PGP

Before starting, it is necessary to distinguish between two noun concepts: PGP (Pretty Good Privacy) and GPG (GNU Privacy Guard). The first one is only a series of protocols and standards, not a specific tool or command. It specifies how to encrypt and decrypt files or contents using specific methods and algorithms. Any tool that implement OpenPGP protocol has the ability of encryption and decription. 
The second one, is an open source free encryption tool under the GNU Project. At present, most Linux distributions have GPG installed by default. To chek the GPG availability, type in a shell:

```
gpg --version
```

or install it manually:

```
sudo apt install gnupg
```

### Generate GPG key pair

GPG uses the concepts of _Public key_ and _Private key_ to implement asymmetric encryption algorithm:

* Public key is used for encryption. This key can be distributed to any organization or idividual and user with this key can only encrypt
* Private key is used for decription, and can only be used to decrypt the information encrypted by the public key paired with the public key.

Before using GPG commands for encryption and decryption, first of all generate the GPG public and private keys:

```
gpg --gen-key
```

And follow the interactive procedure.
This command will generate a pair of RSA key pairs with a length of 4096 and never expire.

After generating the GPG key pair, it is possible to list them by typing:

```
gpg --list-keys
--or--
gpg --list-secret-keys
```

To see a specific key information, it is necessary to pass at least one information between _UserName_, _mailbox_ and _keyID_ (last 16 characters of the key):

```
gpg --list-key "user"
--or--
gpg --list-key "mailbox"
--or--
gpg --list-key <keyID>
```

After generating key pairs, it is possible to use them to encrypt/decrypt files.

### Encypt and decrypt using GPG

To better understand the advantages of SOPs command, let's start by using directly the gpg command to encrypt:

```
gpg -e -r <keyID/user/mailbox> <fileName>
```

Where:

* _-e_: Parameter that specifies the operation to execute (encryption)
* _-r_: Parameter that specifies the public key information to be used for encryption

Once done, a new file in the same directory is created with the extension _<fileName>.gpg_

It contains encrypted content and trying to view it returns a garbed code, which is the formal working mode of GPG.

It can be also used to decrypt an ecrypted file, as follow:

```
gpg -d <fileName>.gpg
```

Where:

* _-d_: Parameter that specify the operation is a decryption. Unlike the encryption step, it is not necessary to specify the private key ID, because GPG writes the private key information required for decryption into the encrypted file at the same time of encryption. Of course you can still use use it by adding the explic parameter _-r <privateKeyID>_.

Once done, the decrypted file is printed to shell.

The GPG command for encryption perform an operation for the whole file, making the encrypted file completely unmaintainable.

### Simple use of SOPs

To encrypt a file using SOPs and GPG generated key, type:

```
sops --encrypt --in-place --pgp <keyID> <fileName>
```

Where:

* _--encrypt, -e_: Parameter that tells sops to encrypt
* _--in-place, -i_: Parameter that specifies to directly replace the encrypted content with the original file content. If this parameter is not specified, the encrypted content will be displayed to screen
* _--pgp, -p_: Parameter that specifies the PGP public key ID to use when encrypting

Opening the encrypted file, it is possible to see that SOPs only encrypts the values of the yaml file and retain all key information, which makes even the encrypted file still retain strong readability.

Additionally, SOPs adds a new key value which is used to save some encrypted information used to encrypt the file.

Of course it is also possible to decrypt the file typing:

```
sops --decrypt <fileName>
```

Where:

* _--decrypt, -d_: Parameter that specifies the decryption operation. It is possible to pass the private key ID using the _-p_ parameter during decryption step, but usually this is not necessary, because like GPG, SOPs automatically records the private key ID required during decryption while encrypting files.

## Using Helm secrets plug-in for deployment

The Helm secrets plug-in creates a new subcommand for the Helm, all the subcommands about the plug-in are called through the helm secrets format.

_Before proceding with these steps, be sure to match the prerequisites, install the helm secrets plugin and check the sops and gpg commands._

To deploy the encrypted application on kubernetes, follow these steps:

* Create encryption key pair:

```
gpg --gen-key
```

And follow the procedure, adding correct information and defining a passphrase. This can be used to decrypt the file manually during deploy

* Encrypt the values.yaml file into the new file secrets.values.yaml:
 
```
sops --encrypt --in-place --pgp <keyID> secrets.values.yaml
```

* Deploy or upgrade the k3s distro:

```
--DEPLOY--
helm secrets install tim-vanilla /home/andrea/tim-edgex/tim-edgex/ -f /home/andrea/tim-edgex/tim-edgex/secrets.values.yaml --kubeconfig /etc/rancher/k3s/k3s.yaml
--UPGRADE--
helm secrets upgrade tim-vanilla /home/andrea/tim-edgex/tim-edgex/ -f /home/andrea/tim-edgex/tim-edgex/secrets.values.yaml --kubeconfig /etc/rancher/k3s/k3s.yaml
```

At this point the passphrase will be required to decrypt the file and proceed with the deployment. If the passphrase is wrong the deployment will fail and the encrypted file should not be decrypted, so be careful not to lose this key

### Related links

For more details about encryption/decryption, take a look to the following link: https://developpaper.com/helm-secrets-of-helm-plug-in-encrypt-your-values-file-with-pgp-idcf/