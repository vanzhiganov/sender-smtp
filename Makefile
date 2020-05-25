.PHONY: build run clean

build:
	go build -o build/pkg-debian/usr/bin/sender-smtp-api

build_deb:
	go build -o build/pkg-debian/usr/bin/sender-smtp-api
	VERSION=`build/pkg-debian/usr/bin/sender-smtp-api version` sed -i -e "s/#Version: .../Version: ${VERSION}/g" build/pkg-debian/DEBIAN/control
	SIZE=`du -k build/pkg-debian/usr/bin/sender-smtp-api | awk '{print $1}'` sed -i -e "s/#Installed-Size: .../Installed-Size: ${SIZE}/g" build/pkg-debian/DEBIAN/control

	cd build/pkg-debian && find . -type f ! -regex '..hg.' ! -regex '..git.' ! -regex '.?debian-binary.' ! -regex '.?DEBIAN..' -printf '%P ' | xargs md5sum > DEBIAN/md5sums
	VERSION=`build/pkg-debian/usr/bin/sender-smtp-api version` dpkg -b build/pkg-debian/ dist/sender-smtp-api_${VERSION}_amd64.deb

run:
	go run main.go struct_config.go struct_request.go struct_response.go config.go endpoint_sendmail_v1.go

clean:
	go clean ./...
