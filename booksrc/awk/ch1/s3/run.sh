#!/bin/bash

awk -f printf_1.awk ../emp.data
awk -f printf_2.awk ../emp.data
awk -f printf_3.awk ../emp.data | sort -n

