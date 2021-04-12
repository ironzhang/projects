#include <stdio.h>

#include "codec.h"
#include "protobuf_c.h"

#pragma pack(1)
typedef struct head_t {
	uint8_t msg_id;
	uint16_t sequence;
}head_t;
#pragma pack()


field_meta head_fields_meta[] = {
	{"msg_id", "uint8", 1},
	{"sequence", "uint16", 2},
};
message_meta head_meta = {
	"Head",
	2,
	head_fields_meta,
};

static uint8_t buf[1024];

int main() {
	head_t head = {8, 255+1};
	int n = message_encode(&head_meta, &head, buf, 1024);
	if (n <= 0) {
		fprintf(stderr, "message encode: %d", n);
		return 0;
	}
	for (int i = 0; i < n; i++) {
		fprintf(stdout, "%d,", buf[i]);
	}
	fprintf(stdout, "\n");
	return 0;
}
