RM = rm
CC = cc

TARGET = test
CFLAGS = -Wall -g
LFLAGS =

all: $(TARGET)

$(TARGET): thingpb.c
	$(CC) $(CFLAGS) -o $@ $^ $(LFLAGS)

clean:
	$(RM) $(TARGET)
