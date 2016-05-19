%if 0%{?rhel} == 7
  %define dist .el7
%endif
%define _unpackaged_files_terminate_build 0
Name: relevy
Version: 0.1
Release:	1%{?dist}
Summary: A golang backend for placing host information in mongodb.

License: GPLv2
URL: https://github.com/Jmainguy/relevy
Source0: relevy.tar.gz
Requires(pre): shadow-utils

%description
A golang backend for slurping information from the host, and placing it in mongodb to be consumed by frontends.

%prep
%setup -q -n relevy
%install
mkdir -p $RPM_BUILD_ROOT/usr/sbin
mkdir -p $RPM_BUILD_ROOT/opt/relevy
mkdir -p $RPM_BUILD_ROOT/usr/lib/systemd/system
mkdir -p $RPM_BUILD_ROOT/etc/relevy
mkdir -p $RPM_BUILD_ROOT/etc/rc.d/init.d/
install -m 0755 $RPM_BUILD_DIR/relevy/relevy %{buildroot}/usr/sbin
install -m 0755 $RPM_BUILD_DIR/relevy/service/relevy.sysv %{buildroot}/etc/rc.d/init.d/relevy
install -m 0644 $RPM_BUILD_DIR/relevy/service/relevy.service %{buildroot}/usr/lib/systemd/system
install -m 0644 $RPM_BUILD_DIR/relevy/config.yaml %{buildroot}/etc/relevy/
install -m 0644 $RPM_BUILD_DIR/relevy/info.yaml %{buildroot}/etc/relevy/

%files
/usr/sbin/relevy
%if 0%{?rhel} == 6
  /etc/rc.d/init.d/relevy
%else
  /usr/lib/systemd/system/relevy.service
%endif
%dir /opt/relevy
%dir /etc/relevy
%config(noreplace) /etc/relevy/config.yaml
%config(noreplace) /etc/relevy/info.yaml

%pre
getent group relevy >/dev/null || groupadd -r relevy
getent passwd relevy >/dev/null || \
    useradd -r -g relevy -d /opt/relevy -s /sbin/nologin \
    -c "User to run relevy service" relevy
exit 0

%post
chown -R relevy:relevy /opt/relevy
if [ -f /usr/bin/systemctl ]; then
  systemctl daemon-reload
fi

%changelog

