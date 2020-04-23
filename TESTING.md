# Setting up Chef Zero

Run Chef:

```
mkdir /tmp/chef
# For persistance
docker run --rm -p 8889:8889 --name chef -v /tmp/chef:/work -it rubygem/chef-zero -H 0.0.0.0 -p 8889 -l debug --file-store /work
# Without persistance
docker run --rm -p 8889:8889 --name chef -it rubygem/chef-zero -H 0.0.0.0 -p 8889 -l debug
```

Setup Environment:

Generate a private key to use in the key material var

```
openssl genrsa -out key.pem 2048
export CHEF_SERVER_URL=http://127.0.0.1:8889/
export CHEF_CLIENT_NAME=test
export CHEF_KEY_MATERIAL="$(cat key.pem)"
```

Run: 

```
make testacc
```
