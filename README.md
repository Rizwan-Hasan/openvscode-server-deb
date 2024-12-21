# openvscode-server-deb

### Download openvscode-server

```bash
wget -c "https://github.com/gitpod-io/openvscode-server/releases/download/openvscode-server-v1.96.0/openvscode-server-v1.96.0-linux-arm64.tar.gz"
```

### Extract archive

```shell
tar xzvf openvscode-server-v1.96.0-linux-arm64.tar.gz
```

```shell
mv -vf openvscode-server-v1.96.0-linux-arm64 ./pkgroot/opt/openvscode-server
```

### Get License

```shell
wget -c "https://raw.githubusercontent.com/gitpod-io/openvscode-server/refs/heads/main/LICENSE.txt"
```

```shell
mv -vf LICENSE.txt ./pkgroot/usr/share/licenses/openvscode-server.txt
```

### Add DEBIAN folder

```shell
cp -avrf debian-files ./pkgroot/DEBIAN
```

### Update architecture

```shell
sed -i 's\ARCHITECTURE\amd64\g' ./pkgroot/DEBIAN/control  # For Intel/AMD
sed -i 's\ARCHITECTURE\arm64\g' ./pkgroot/DEBIAN/control  # For ARM64
```

### Update version

```shell
sed -i 's\VERSION\1.96.0\g' ./pkgroot/DEBIAN/control
```

### Fix permissions

```shell
chmod 755 -R pkgroot
```

### Build deb

```shell
dpkg-deb --build pkgroot openvscode-server-v1.96.0-arm64.deb
```