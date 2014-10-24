sysd
====

## Build sysd with docker

```
git clone https://github.com/hacking-thursday/sysd
cd sysd
sudo docker build -t sysd .
```

## Get sysd from docker image

```
sudo docker run -v "$PWD:/dist" sysd cp /usr/local/bin/sysd /dist
```
