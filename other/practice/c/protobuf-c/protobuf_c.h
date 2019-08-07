#ifndef PROTOBUF_C_H 
#define PROTOBUF_C_H

#include <stdint.h>
#include <stddef.h>

typedef struct field_meta {
	const char *name;
	const char *type;
	int tag;
}field_meta;

typedef struct message_meta {
	const char *name;
	int field_count;
	field_meta *fields;
}message_meta;

int message_encode(message_meta *meta, void *msg, uint8_t *buf, size_t sz);

int message_decode(message_meta *meta, void *msg, uint8_t *buf, size_t sz);

#endif
