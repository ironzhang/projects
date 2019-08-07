#ifndef LEXER_H
#define LEXER_H

#include "lex_comm.h"

enum E_TOKEN_TYPE{
	E_TOKEN_TYPE_EOF,
	E_TOKEN_TYPE_ERR,
	E_TOKEN_TYPE_ID,
	E_TOKEN_TYPE_NUM,
	E_TOKEN_TYPE_LE,
	E_TOKEN_TYPE_NE,
	E_TOKEN_TYPE_LT,
	E_TOKEN_TYPE_EQ,
	E_TOKEN_TYPE_GE,
	E_TOKEN_TYPE_GT,
};

int init_lexer();

int next_token(token_t *token);

#endif

