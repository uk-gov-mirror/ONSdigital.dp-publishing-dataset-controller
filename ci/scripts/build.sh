#!/bin/bash -eux

pushd dp-publishing-dataset-controller
  make build
  cp build/dp-publishing-dataset-controller Dockerfile.concourse ../build
popd
