#include "lexer.h"
#include "lex_move.h"
#include "lex_state.h"
#include <stdio.h>
#include <assert.h>
#include <string.h>

static int le_token_func(const char *str, token_t *token)
{
	token->type = E_TOKEN_TYPE_LE;
	token->str = strdup(str);
	return 0;
}

static int ne_token_func(const char *str, token_t *token)
{
	token->type = E_TOKEN_TYPE_NE;
	token->str = strdup(str);
	return 0;
}

static int lt_token_func(const char *str, token_t *token)
{
	token->type = E_TOKEN_TYPE_LT;
	token->str = strdup(str);
	return 0;
}

static int eq_token_func(const char *str, token_t *token)
{
	token->type = E_TOKEN_TYPE_EQ;
	token->str = strdup(str);
	return 0;
}

static int ge_token_func(const char *str, token_t *token)
{
	token->type = E_TOKEN_TYPE_GE;
	token->str = strdup(str);
	return 0;
}

static int gt_token_func(const char *str, token_t *token)
{
	token->type = E_TOKEN_TYPE_GT;
	token->str = strdup(str);
	return 0;
}

static int id_token_func(const char *str, token_t *token)
{
	token->type = E_TOKEN_TYPE_ID;
	token->str = strdup(str);
	return 0;
}

static int num_token_func(const char *str, token_t *token)
{
	token->type = E_TOKEN_TYPE_NUM;
	token->str = strdup(str);
	return 0;
}

static int err_token_func(const char *str, token_t *token)
{
	token->type = E_TOKEN_TYPE_ERR;
	token->str = strdup(str);
	return 0;
}

static void register_lexer_states()
{
	register_lex_state(2, 0, le_token_func);
	register_lex_state(3, 0, ne_token_func);
	register_lex_state(4, 1, lt_token_func);
	register_lex_state(5, 0, eq_token_func);
	register_lex_state(7, 0, ge_token_func);
	register_lex_state(8, 1, gt_token_func);

	register_lex_state(11, 1, id_token_func);

	register_lex_state(19, 1, num_token_func);
	register_lex_state(20, 1, num_token_func);
	register_lex_state(21, 1, num_token_func);

	register_lex_state(22, 0, err_token_func);
}

static void set_lexer_move_edges()
{
	lex_move_set_edge(0, '<', 1);
	lex_move_set_edge(0, '=', 5);
	lex_move_set_edge(0, '>', 6);
	lex_move_set_letter_edges(0, 10);
	lex_move_set_digit_edges(0, 13);
	lex_move_set_other_edges(0, 22);

	lex_move_set_edge(1, '=', 2);
	lex_move_set_edge(1, '>', 3);
	lex_move_set_other_edges(1, 4);

	lex_move_set_edge(6, '=', 7);
	lex_move_set_other_edges(6, 8);

	lex_move_set_letter_edges(10, 10);
	lex_move_set_digit_edges(10, 10);
	lex_move_set_other_edges(10, 11);

	lex_move_set_edge(13, '.', 14);
	lex_move_set_edge(13, 'E', 16);
	lex_move_set_digit_edges(13, 13);
	lex_move_set_other_edges(13, 20);

	lex_move_set_digit_edges(14, 15);
	lex_move_set_other_edges(14, 22);

	lex_move_set_edge(15, 'E', 16);
	lex_move_set_digit_edges(15, 15);
	lex_move_set_other_edges(15, 21);

	lex_move_set_edge(16, '+', 17);
	lex_move_set_edge(16, '-', 17);
	lex_move_set_digit_edges(16, 18);
	lex_move_set_other_edges(16, 22);

	lex_move_set_digit_edges(17, 18);
	lex_move_set_other_edges(17, 22);

	lex_move_set_digit_edges(18, 18);
	lex_move_set_other_edges(18, 19);
}

int init_lexer()
{
	register_lexer_states();
	set_lexer_move_edges();
	return 0;
}

int next_token(token_t *token)
{
	char ch;
	int i, state_id;
	char token_str[256];
	lex_state_t *state;

	assert(token != NULL);

	do
		ch = getc(stdin);
	while(ch == ' ' || ch == '\t' || ch == '\n');

	if(ch == EOF){
		token->type = E_TOKEN_TYPE_EOF;
		token->str = "";
		return 0;
	}

	i = 0;
	state_id = 0;
	while(1){
		token_str[i++] = ch;
		state_id = lex_move_next_state(state_id, ch);
		state = get_lex_state(state_id);
		if(state->is_fini){
			if(state->need_back){
				i--;
				ungetc(ch, stdin);
			}
			token_str[i] = '\0';
			return state->fn_op(token_str, token);
		}
		ch = getc(stdin);
	}

	return 0;
}

