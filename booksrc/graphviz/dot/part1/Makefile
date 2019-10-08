# Makefile

all_dot = $(wildcard *.dot)
all_png = $(all_dot:.dot=.png)
all_ps  = $(all_dot:.dot=.ps)

all: $(all_png) $(all_ps)

%.png: %.dot
	dot -Tpng $^ -o $@

%.ps: %.dot
	dot -Tps $^ -o $@

clean:
	rm -rf $(all_png) $(all_ps)

