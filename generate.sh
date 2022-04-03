#!/bin/bash -ex
docker run -v $PWD:/mnt stefansundin/ppastats ppastats stefansundin aws-rotate-key -o /mnt
