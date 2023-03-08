#include "codec.h"

int encode_uint8(uint8_t x, uint8_t* buf, size_t sz) {
	if (sz < sizeof(uint8_t))
		return -1;
	buf[0] = x;
	return sizeof(uint8_t);
}

int decode_uint8(uint8_t *x, uint8_t* buf, size_t sz) {
	if (sz < sizeof(uint8_t))
		return -1;
	*x = buf[0];
	return sizeof(uint8_t);
}

int encode_uint16(uint16_t x, uint8_t* buf, size_t sz) {
	if (sz < sizeof(uint16_t))
		return -1;
	buf[0] = (uint8_t)(x);
	buf[1] = (uint8_t)(x >> 8);
	return sizeof(uint16_t);
}

int decode_uint16(uint16_t *x, uint8_t* buf, size_t sz) {
	if (sz < sizeof(uint16_t))
		return -1;
	*x = (uint16_t)(buf[0]);
	*x |= (uint16_t)(buf[1]) << 8;
	return sizeof(uint16_t);
}
