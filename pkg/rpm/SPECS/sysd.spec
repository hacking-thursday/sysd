Summary: the daemon who supplies firsthand system data
Name: sysd
Version: 0.5.0
Release: 6%{?dist}
License: Apache License, Version 2.0
Group: Applications/System
Source0: http://github.com/hacking-thursday/sysd/archive/sysd-0.5.0.tar.gz
#Source1: 
Url: https://github.com/hacking-thursday/sysd
#Patch: 
Buildroot: %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)
BuildRequires: golang
#Requires(preun): 
#Requires(post): 

%description
sysd implements a light dependeny daemon in golang, and provides /proc,/sys the firsthand system data in json/xml/... common formats with a high-level RESTful HTTP API. With sysd, application and plugin developers are able to save their works from parsing variant output from low-level unix command tools, and dependencies.

%prep
%setup -q

%build
%configure

make %{?_smp_mflags}

%install
rm -rf $RPM_BUILD_ROOT

make install DESTDIR=$RPM_BUILD_ROOT

%post

%preun

%clean
rm -rf $RPM_BUILD_ROOT

%files
%defattr(-,root,root)
%doc LICENSE
#%attr(0644,root,root) %{_sysconfdir}/init.d/sysd
#%{_bindir}/*
%{_sbindir}/*

%changelog
* Sun Nov 09 2014 Chun-Yu Lee (Mat) <matlinuxer2@gmail.com>
- Initial release

