#include "lex_move.h"
#include "lex_comm.h"
#include <assert.h>

static int gs_move_other_edge[MAX_LEX_STATE_NUM];
static int gs_move_table[MAX_LEX_STATE_NUM][256];

void lex_move_set_edge(int src_state, int ch, int dst_state)
{
	assert(ch >= 0 && ch < 256);
	assert(src_state >= 0 && src_state < MAX_LEX_STATE_NUM);
	assert(dst_state >= 0 && dst_state < MAX_LEX_STATE_NUM);
	gs_move_table[src_state][ch] = dst_state;
}

void lex_move_set_letter_edges(int src_state, int dst_state)
{
	int i;
	assert(src_state >= 0 && src_state < MAX_LEX_STATE_NUM);
	assert(dst_state >= 0 && dst_state < MAX_LEX_STATE_NUM);
	for(i = 'A'; i <= 'Z'; i++)
		gs_move_table[src_state][i] = dst_state;
	for(i = 'a'; i <= 'z'; i++)
		gs_move_table[src_state][i] = dst_state;
}

void lex_move_set_digit_edges(int src_state, int dst_state)
{
	int i;
	assert(src_state >= 0 && src_state < MAX_LEX_STATE_NUM);
	assert(dst_state >= 0 && dst_state < MAX_LEX_STATE_NUM);
	for(i = '0'; i <= '9'; i++)
		gs_move_table[src_state][i] = dst_state;
}

void lex_move_set_other_edges(int src_state, int dst_state)
{
	assert(src_state >= 0 && src_state < MAX_LEX_STATE_NUM);
	assert(dst_state >= 0 && dst_state < MAX_LEX_STATE_NUM);
	gs_move_other_edge[src_state] = dst_state;
	/*
	for(i = 0; i < 256; i++)
		if(gs_move_table[src_state][i] == 0)
			gs_move_table[src_state][i] = dst_state;
	*/
}

int lex_move_next_state(int state, int ch)
{
	int next_state = 0;
	assert(state >= 0 && state < MAX_LEX_STATE_NUM);
	if(ch >= 0 && ch < 256)
		next_state = gs_move_table[state][ch];
	if(next_state == 0)
		next_state = gs_move_other_edge[state];
	return next_state;
}

