#!/bin/bash
# Created by:Stitch-Zhang
# @https://gitee.com/Stitchtor/ptm
version=0.1.3
URL=""
URL_BASE=https://gitee.com/Stitchtor/ptm/attach_files/681282/download/

URLA=${URL_BASE}ptm_v${version}-Linux_x86_64.tar.gz
URLB=${URL_BASE}ptm_v${version}-Linux_i386.tar.gz

SYSTYPE=$(echo $(uname -m) | grep "64")

if [ SYSTYPE ]
then
	URL=${URLA}
else
	URL=${URLB}
fi

PWD=$(pwd)/tmp

mkdir ${PWD} && cd ${PWD}
curl -L ${URL} -o ptm${version}.tar.gz && tar xvf ptm${version}.tar.gz && cp ptm /usr/bin/
rm -rf ${PWD}
echo "Installed"
echo "Type ptm to run"