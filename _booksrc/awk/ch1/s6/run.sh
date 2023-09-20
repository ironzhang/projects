#!/bin/bash

awk -f ./c1.awk ../emp.data
awk -f ./interest1.awk ./interest.data
awk -f ./interest2.awk ./interest.data
awk -f ./reverse1.awk ../emp.data
awk -f ./reverse2.awk ../emp.data
