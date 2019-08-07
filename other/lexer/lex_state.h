#ifndef LEX_STATE_H
#define LEX_STATE_H

#include "lex_comm.h"

typedef int (*fn_lex_state_op_t)(const char *str, token_t *token);

typedef struct lex_state_t{
	int state_id;
	char is_fini;
	char need_back;
	fn_lex_state_op_t fn_op;
}lex_state_t;

void register_lex_state(int state_id, char need_back, fn_lex_state_op_t fn_op);

lex_state_t *get_lex_state(int state_id);

#endif

