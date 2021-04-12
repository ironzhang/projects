RM = rm
CC = g++

TARGET = test
CFLAGS = -Wall -g
LFLAGS =

all: $(TARGET)

$(TARGET): main.c protobuf_c.c codec.c
	$(CC) $(CFLAGS) -o $@ $^ $(LFLAGS)

clean:
	$(RM) $(TARGET)
