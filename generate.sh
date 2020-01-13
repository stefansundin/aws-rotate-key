#!/bin/bash -ex
multipass exec ppastats -- ppastats stefansundin aws-rotate-key -o "$(pwd)"
