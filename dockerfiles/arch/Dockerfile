FROM base/archlinux

RUN pacman -Sy --noconfirm go git libunistring patch make

RUN curl -O https://aur.archlinux.org/packages/sy/sysd/PKGBUILD

RUN makepkg -ci --asroot --noconfirm

EXPOSE 8

ENTRYPOINT ["/usr/sbin/sysd"]
