#!/bin/sh

gometalinter \
--deadline=90s \
--concurrency=6 \
--vendor \
--cyclo-over=15 \
--tests \
--exclude="exported (var)|(method)|(const)|(type)|(function) [A-Za-z\.0-9]* should have comment" \
--exclude="Id.* should be .*ID" \
./...
