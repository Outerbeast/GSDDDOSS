#!/bin/bash
cd src
go build -ldflags="-s -w" -o GSDDOSS .
mv GSDDOSS ../GSDDOSS
