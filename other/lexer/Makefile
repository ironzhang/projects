TARGET = lexer

all: $(TARGET)

$(TARGET): lex_state.c lex_move.c lexer.c main.c
	gcc -g -Wall -o $@ $^

clean:
	rm $(TARGET)

