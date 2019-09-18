#!/bin/bash

awk '$3 > 0 { print $1, $2 * $3 }' emp.data

awk '$3 == 0 { print $1 }' emp.data

awk '$3 == 0' emp.data

awk '{ print $1 }' emp.data
