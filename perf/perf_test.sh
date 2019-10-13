#!/usr/bin/env bash

timeout 60 siege -b -c 100 http://localhost:5000/
