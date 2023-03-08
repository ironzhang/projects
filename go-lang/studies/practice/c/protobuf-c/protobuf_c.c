#include "protobuf_c.h"

#include <string.h>

#include "codec.h"

message_meta metas[] = {
	{"int8"},
	{"uint8"},
	{"int16"},
	{"uint16"},
};

static int message_encode_uint8(void *m, uint8_t *buf, size_t sz) {
	uint8_t *p = (uint8_t *)m;
	return encode_uint8(*p, buf, sz);
}

static int message_encode_uint16(void *m, uint8_t *buf, size_t sz) {
	uint16_t *p = (uint16_t *)m;
	return encode_uint16(*p, buf, sz);
}

static int field_encode(field_meta *meta, void *msg, uint8_t *buf, size_t sz) {
	for (size_t i = 0; i < sizeof(metas)/sizeof(metas[0]); i++) {
		if (strcmp(meta->type, metas[i].name) == 0)
			return message_encode(&metas[i], msg, buf, sz);
	}
	return -1;
}

int message_encode(message_meta *meta, void *msg, uint8_t *buf, size_t sz) {
	if ((strcmp(meta->name, "int8") == 0) || (strcmp(meta->name, "uint8") == 0))
		return message_encode_uint8(msg, buf, sz);
	else if ((strcmp(meta->name, "int16") == 0) || (strcmp(meta->name, "uint16") == 0))
		return message_encode_uint16(msg, buf, sz);

	int sum = 0;
	uint8_t *ptr = (uint8_t *)msg;
	for (int i = 0; i < meta->field_count; i++) {
		int n = field_encode(&meta->fields[i], ptr, buf, sz);
		if (n < 0) {
			return n;
		}
		sum += n;
		ptr += n;
		buf += n;
		sz -= n;
	}
	return sum;
}

//static int message_encode_int32(void *msg, buffer *buf) {
//	int32_t *p = (int32_t *)(msg);
//	encode_int32(buf->ptr + buf->offset, *p);
//	return sizeof(int32_t)
//}

//int message_encode(message_prototype *proto, void *msg, buffer *buf) {
//	int size = 0;
//	int offset = 0;
//	unsigned char *p = (unsigned char *)msg;
//
//	for (int i = 0; i < proto->filed_count; i++) {
//		size = field_prototype_encode(&proto->fileds[i], p+offset, buf);
//		offset += size;
//	}
//}
//
//int field_prototype_encode(field_prototype *proto, void *msg, buffer *buf) {
//	message_prototype *message_proto = find_message_proto_type(proto->type);
//	return message_encode(message_prototype, msg, buf)
//}
