#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-publishing-dataset-controller
  make audit
popd