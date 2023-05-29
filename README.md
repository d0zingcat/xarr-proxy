# XarrProxy
A proxy layer for sonarr or something.

Inspired by [Jproxy](https://github.com/LuckyPuppy514/jproxy).

The differences between XarrProxy and Jproxy:
- XarrProxy uses golang while Jproxy uses java.
- XarrProxy uses toml to manage configs while Jproxy uses JSON.
- XarrProxy uses bcrypto to handle login/password encryption while Jproxy uses md5.

Also thanks for the following projects, some ideas and designs are borrowed from there:

[Xarr-Rss](https://github.com/xiaoyi510/xarr-rss)
