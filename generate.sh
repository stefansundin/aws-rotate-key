#!/bin/bash -ex
docker run -v $PWD:/mnt stefansundin/ppastats stefansundin aws-rotate-key -o /mnt
