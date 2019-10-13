#!/usr/bin/env bash

timeout 30 siege -b -c 100 http://localhost:5000/
