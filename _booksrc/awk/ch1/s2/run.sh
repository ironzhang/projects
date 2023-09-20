#!/bin/bash

awk -f print_line.awk emp.data
awk -f print_field.awk emp.data
awk -f print_nf.awk emp.data
awk -f print_nr.awk emp.data
awk -f print_text.awk emp.data
