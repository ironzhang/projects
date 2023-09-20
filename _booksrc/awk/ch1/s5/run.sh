#!/bin/bash

awk -f ./calc_1.awk ../emp.data
awk -f ./calc_2.awk ../emp.data
awk -f ./calc_3.awk ../emp.data
awk -f ./calc_4.awk ../emp.data
awk -f ./names.awk ../emp.data
awk -f ./last_line.awk ../emp.data
awk -f ./name_len.awk ../emp.data
awk -f ./wc.awk ../emp.data
