#!/bin/bash

awk -f select_1.awk ../emp.data
awk -f select_2.awk ../emp.data
awk -f select_3.awk ../emp.data
awk -f select_4.awk ../emp.data
awk -f select_5.awk ../emp.data
awk -f select_6.awk ../emp.data
awk -f select_7.awk ../emp.data
awk -f validate.awk ../emp.data
awk -f title.awk ../emp.data
