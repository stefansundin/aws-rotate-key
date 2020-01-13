#!/bin/bash -ex
multipass launch -n ppastats --cloud-init cloud-init.yaml xenial
multipass mount . ppastats
