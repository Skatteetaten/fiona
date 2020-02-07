#!/usr/bin/env bash

make clean
make test-xml
make test-coverage
go test -short -coverprofile=bin/cov.out `go list ./... | grep -v vendor/`
make

#Create archive
ARCHIVEFOLDER="fiona_app"
mkdir -p bin/$ARCHIVEFOLDER
cp bin/fiona bin/$ARCHIVEFOLDER/fiona
#Copy metadata if available
if test -f "metadata/openshift.json"; then
  mkdir bin/$ARCHIVEFOLDER/metadata
  cp metadata/openshift.json bin/$ARCHIVEFOLDER/metadata/openshift.json
else
  echo "NB: Openshift metadata not found"
fi
echo "Creating archive fiona.zip:"
cd bin
zip -r fiona.zip $ARCHIVEFOLDER
#cleanup
rm -rf $ARCHIVEFOLDER