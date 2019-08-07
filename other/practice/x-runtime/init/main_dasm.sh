#!/bin/bash

go tool objdump -s "main\.init" init | less
