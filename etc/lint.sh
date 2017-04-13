#!/bin/sh

gometalinter --deadline 10s $1 \
   | grep -v ID \
   | grep -v comment \
   | grep -v JSON

