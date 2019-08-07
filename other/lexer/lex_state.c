#include "lex_state.h"
#include <assert.h>
#include <stddef.h>

static lex_state_t gs_lex_states[MAX_LEX_STATE_NUM];

void register_lex_state(int state_id, char need_back, fn_lex_state_op_t fn_op)
{
	assert(fn_op != NULL);
	assert(state_id >= 0 && state_id < MAX_LEX_STATE_NUM);
	gs_lex_states[state_id].state_id = state_id;
	gs_lex_states[state_id].is_fini = 1;
	gs_lex_states[state_id].need_back = need_back;
	gs_lex_states[state_id].fn_op = fn_op;
}

lex_state_t *get_lex_state(int state_id)
{
	assert(state_id >= 0 && state_id < MAX_LEX_STATE_NUM);
	return &gs_lex_states[state_id];
}

