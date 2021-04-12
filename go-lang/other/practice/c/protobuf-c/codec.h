#ifndef CODEC_H
#define CODEC_H

#include <stdint.h>
#include <stddef.h>

int encode_uint8(uint8_t x, uint8_t* buf, size_t sz);
int decode_uint8(uint8_t *x, uint8_t* buf, size_t sz);

int encode_uint16(uint16_t x, uint8_t* buf, size_t sz);
int decode_uint16(uint16_t *x, uint8_t* buf, size_t sz);

#endif
