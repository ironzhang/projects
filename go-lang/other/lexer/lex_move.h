#ifndef LEX_MOVE_H
#define LEX_MOVE_H

void lex_move_set_edge(int src_state, int ch, int dst_state);
void lex_move_set_letter_edges(int src_state, int dst_state);
void lex_move_set_digit_edges(int src_state, int dst_state);
void lex_move_set_other_edges(int src_state, int dst_state);

int lex_move_next_state(int state, int ch);

#endif

