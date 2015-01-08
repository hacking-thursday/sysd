FROM fedora:20

RUN yum install -y wget

RUN wget -O /etc/yum.repos.d/home:matlinuxer2.repo http://download.opensuse.org/repositories/home:matlinuxer2/Fedora_20/home:matlinuxer2.repo

RUN yum install -y sysd

RUN chmod u+x /usr/sbin/sysd

EXPOSE 8

ENTRYPOINT ["/usr/sbin/sysd"]
