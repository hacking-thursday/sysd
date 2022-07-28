build:
	( install -d SOURCES/; cp ../tgz/*.tar.gz -v SOURCES/ )
	rpmbuild --define "_topdir `pwd`" --quiet --clean -ba SPECS/sysd.spec

setup:
	yum groupinstall "Development Tools"
	yum install epel-release
	yum install golang
