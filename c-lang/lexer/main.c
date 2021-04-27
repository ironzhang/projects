#include <stdio.h>
#include <assert.h>
#include "lexer.h"

const char *get_token_type_str(int token_type)
{
	const char *str = "UNKNOWN";
	switch(token_type){
		case E_TOKEN_TYPE_EOF:
			str = "E_TOKEN_TYPE_EOF";
			break;
		case E_TOKEN_TYPE_ERR:
			str = "E_TOKEN_TYPE_ERR";
			break;
		case E_TOKEN_TYPE_ID:
			str = "E_TOKEN_TYPE_ID";
			break;
		case E_TOKEN_TYPE_NUM:
			str = "E_TOKEN_TYPE_NUM";
			break;
		case E_TOKEN_TYPE_LE:
			str = "E_TOKEN_TYPE_LE";
			break;
		case E_TOKEN_TYPE_NE:
			str = "E_TOKEN_TYPE_NE";
			break;
		case E_TOKEN_TYPE_LT:
			str = "E_TOKEN_TYPE_LT";
			break;
		case E_TOKEN_TYPE_EQ:
			str = "E_TOKEN_TYPE_EQ";
			break;
		case E_TOKEN_TYPE_GE:
			str = "E_TOKEN_TYPE_GE";
			break;
		case E_TOKEN_TYPE_GT:
			str = "E_TOKEN_TYPE_GT";
			break;
		default:
			break;
	}
	return str;
}

int main()
{
	int ret;
	token_t token;

	init_lexer();

	do{
		ret = next_token(&token);
		if(ret != 0){
			printf("get next token error, type[%d], str[%s]", token.type, token.str);
			continue;
		}

		if(token.type == E_TOKEN_TYPE_EOF)
			break;

		printf("token: type[%d %s], str[%s]\n", token.type, get_token_type_str(token.type), token.str);

	}while(1);

	return 0;
}

