#!/bin/sh

# mcast - Command line tool and library for testing multicast traffic
# flows and stress testing networks and devices.
# Copyright (C) 2018 Will Smith
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

BINDIR=binaries
ARCH=amd64
CLI=mcast

mkdir -p $BINDIR
rm -f $BINDIR/mac/*
rm -f $BINDIR/linux/*
rm -f $BINDIR/windows/*

env GOOS=linux GOARCH=amd64 go build -o $BINDIR/linux/$CLI
env GOOS=darwin GOARCH=amd64 go build -o $BINDIR/mac/$CLI
env GOOS=windows GOARCH=amd64 go build -o $BINDIR/windows/$CLI.exe
