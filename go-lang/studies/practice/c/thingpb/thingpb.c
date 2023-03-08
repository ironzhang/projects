#include <stdio.h>

#include "protobuf_c.h"

message_meta head_message = {
	"Head",
	2,
	{
		{"MsgName", "string", 1},
		{"Sequence", "int32", 2},
	}
};

message_meta error_notice_message = {
	"ErrorNotice",
	3,
	{
		{"Code", "int32", 1},
		{"Details", "string", 2},
		{"Desc", "string", 3},
	}
};

void print_message_meta(message_meta *m) {
	printf("%s\n", m->name);
	for (int i = 0; i < m->field_count; i++) {
		printf("\t%s %s %d\n", m->fields[i].name, m->fields[i].type, m->fields[i].tag);
	}
}

int main() {
	print_message_meta(&head_message);
	print_message_meta(&error_notice_message);
	return 0;
}
