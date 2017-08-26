#! /bin/bash

# Copyright 2015 Google Inc. All rights reserved.
# Use of this source code is governed by the Apache 2.0
# license that can be found in the LICENSE file.

set -ex

#if [ ! $(dirname $0) = "." ]; then
 # echo "Must run $(basename $0) from the gce_deployment directory."
  #exit 1
#fi

#if [ -z "$BOOKSHELF_DEPLOY_LOCATION" ]; then
 # echo "Must set \$BOOKSHELF_DEPLOY_LOCATION. For example: BOOKSHELF_DEPLOY_LOCATION=gs://my-bucket/bookshelf-VERSION.tar"
  #exit 1
#fi

#TMP=$(mktemp -d -t gce-deploy-XXXXXX)

# [START cross_compile]
# Cross compile the app for linux/amd64
# GOOS=linux GOARCH=amd64 go build -v -o $TMP/app ./ft-frontend
GOOS=linux GOARCH=amd64 go build -v -o ./app ./ft-frontend
# [END cross_compile]

# [START tar]
# Add the app binary
# tar -c -f $TMP/bundle.tar -C $TMP app
tar -c -f ./bundle.tar -C . app

# Add static files.
#tar -u -f $TMP/bundle.tar -C ./ft-frontend templates
#tar -u -f $TMP/bundle.tar -C ./ft-frontend script
tar -u -f ./bundle.tar -C ./ft-frontend templates
tar -u -f ./bundle.tar -C ./ft-frontend script

# [END tar]

# [START gcs_push]
# BOOKSHELF_DEPLOY_LOCATION is something like "gs://my-bucket/bookshelf-VERSION.tar".
# gsutil cp /bundle.tar gs:://igdsanglah/igdsanglah/ft-frontend.tar
gsutil cp /bundle.tar gs://igdsanglah/igdsanglah/bundle.tar
# [END gcs_push]

rm -rf $TMP