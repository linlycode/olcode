#!/bin/bash
CURDIR=$(pwd)
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
ROOTDIR=$DIR/../
cd $ROOTDIR/web
yarn install --frozen-lockfile
cd $ROOTDIR/devops
python dev.py -t run -e prod -s gw -p ~/work
python dev.py -t deploy -e prod -s web -p ~/work
cd $CURDIR
