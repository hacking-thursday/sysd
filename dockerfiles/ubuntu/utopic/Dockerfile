FROM ubuntu:14.10

RUN echo "deb http://ppa.launchpad.net/matlinuxer2/sysd/ubuntu utopic main" >> /etc/apt/sources.list
RUN echo "deb-src http://ppa.launchpad.net/matlinuxer2/sysd/ubuntu utopic main" >> /etc/apt/sources.list

RUN apt-get update
RUN apt-get install -y --force-yes sysd

EXPOSE 8080

ENTRYPOINT ["/usr/sbin/sysd"]
