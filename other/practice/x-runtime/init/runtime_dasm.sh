#!/bin/bash

go tool objdump -s "runtime\.init" init | less
