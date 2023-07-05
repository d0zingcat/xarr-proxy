# XarrProxy
A proxy layer for sonarr or something.

Inspired by [Jproxy](https://github.com/LuckyPuppy514/jproxy).

## Why invent wheel?
- Not quite familiar with Java/Spring yet
- Java consumes too much memory, epecially for homelab
- Just for fun

## The differences between XarrProxy and Jproxy:
- XarrProxy partially supports the frontend UI, as I only implemented qbittorrent, sonarr and prowlarr.
- XarrProxy is written in golang while Jproxy uses java.
- XarrProxy uses toml to manage configs while Jproxy uses JSON.
- XarrProxy uses bcrypto to handle login/password encryption while Jproxy uses md5.
- XarrProxy now does not use memory cache.

Also thanks for the following projects, some ideas and designs are borrowed from there:

[Xarr-Rss](https://github.com/xiaoyi510/xarr-rss)
