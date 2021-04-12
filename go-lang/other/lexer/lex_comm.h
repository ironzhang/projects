#ifndef LEX_COMM_H
#define LEX_COMM_H

#define MAX_LEX_STATE_NUM 50

typedef struct token_t{
	int type;
	const char *str;
}token_t;

#endif

