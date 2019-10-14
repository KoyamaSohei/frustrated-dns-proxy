# frustrated-dns-proxy


## How to know the value of DNS resolver

- execute below cmd.

```bash
go get github.com/KoyamaSohei/frustrated-dns-proxy
sudo frustrated-dns-proxy
```

- change /etc/resolv.conf (take a backup!)

```
nameserver 127.0.0.1
```

- open browser and visit some websites!