#!/bin/bash
./killNetwork.sh
./setNetwork.sh
node invoke.js
node query.js
