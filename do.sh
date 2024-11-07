#!/bin/sh

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <filename>"
  exit 1
fi

cd /tmp

groff -S -k -etpg -F "/home/runner/.local/groff/font" $1 > `basename $1 .roff`.ps
ps2pdf `basename $1 .roff`.ps
