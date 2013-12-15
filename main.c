#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <math.h>
#include <libgen.h>
#include <ctype.h>

#define RSC7 0x52534337
#define BASE_SIZE 0x2000

uint16_t get_i16_big(char* buffer) {
	uint16_t* val_big = (uint16_t*)buffer;
	return ((*val_big) << 8) | ((*val_big) >> 8);
}

uint32_t get_i32_big(char* buffer) {
	uint16_t lo,hi;
	lo = get_i16_big(buffer);
	hi = get_i16_big(buffer+2);
	return (lo << 16) | (hi);
}

float get_f32(char* buffer) {
	uint32_t i = get_i32_big(buffer);
	float* f = (float*)&i;
	if (isnan(*f)) {
		return 0.0f;
	}
	return *f;
}

float get_f16(char* buffer) {
	uint16_t i = get_i16_big(buffer);
	float f = 0;

	/* Lovingly borrowed from http://stackoverflow.com/a/15118210 */
	uint32_t t1 = i & 0x7fff;
	uint32_t t2 = i & 0x8000;
	uint32_t t3 = i & 0x7c00;
	t1 <<= 13;
	t2 <<= 16;
	t1 += 0x38000000;
	t1 = (t3 == 0 ? 0 : t1);
	t1 |= t2;
	f = *((float*)&t1);
	return f;
}

uint32_t get_part_size(uint32_t flags) {
	uint32_t new_base_size = BASE_SIZE << (int)(flags & 0xf);
	int size = (int)((((flags >> 17) & 0x7f) + (((flags >> 11) & 0x3f) << 1) + (((flags >> 7) & 0xf) << 2) + (((flags >> 5) & 0x3) << 3) + (((flags >> 4) & 0x1) << 4)) * new_base_size);
	for (int i = 0; i < 4; ++i) {
		size += (((flags >> (24 + i)) & 1) == 1) ? (new_base_size >> (1 + i)) : 0;
	}
	return size;
}

uint32_t gfx_ofs;

#define _ADDR(ofs) ((ofs & 0xFFFFFFF) + 0x10)
#define G_ADDR(ofs) ((ofs & 0xFFFFFFF) + gfx_ofs + 0x10)

uint32_t xdd_get_next_drawable(char* xdr_buf);
void dump_drawable(FILE* model_fd, char* xdr_buf, uint32_t drawable_addr, char* model_basename);

int main(int argc, char** argv) {
	if (argc < 2) {
		fprintf(stderr, "usage: xdr2obj input.[xdr,xdd,xft]\n");
		return 1;
	}

	char* xdr_file = argv[1];

	char extension[8];
	strcpy(extension, &xdr_file[strlen(xdr_file)-3]);
	for (char* c = extension; *c; ++c) *c = tolower(*c);

	char xdr_name_cpy[256]; // copy so basename doesn't mangle it
	strcpy(xdr_name_cpy, xdr_file);
	char* model_basename = basename(xdr_name_cpy);
	model_basename[strlen(model_basename)-4] = '\0'; // remove the extension

	/* buffer the xdr file */
	FILE* xdr_fd = fopen(xdr_file, "rb");
	if (!xdr_fd) {
		fprintf(stderr, "unable to open input file: %s\n", xdr_file);
		return 1;
	}

	fseek(xdr_fd, 0, SEEK_END);
	size_t len = ftell(xdr_fd);
	fseek(xdr_fd, 0, SEEK_SET);

	char* xdr_buf = (char*) malloc(len);
	int to_read = len;
	while (to_read > 0) {
		int ret = fread(&xdr_buf[len-to_read], 1, to_read, xdr_fd);
		if (!ret) {
			perror("unable to buffer input file");
			fclose(xdr_fd);
			return 1;
		}
		to_read -= ret;
	}
	fclose(xdr_fd);

	/* check for a valid xdr file */
	uint32_t magic = get_i32_big(&xdr_buf[0]);
	if (magic != RSC7) {
		printf("magic mismatch %x expected %x\n", magic, RSC7);
	}
	uint32_t sys_flags = get_i32_big(&xdr_buf[4*2]);
	uint32_t gfx_flags = get_i32_big(&xdr_buf[4*3]);
	gfx_ofs = get_part_size(sys_flags);

	char model_name[256];
	sprintf(model_name, "%s.obj", model_basename);
	FILE* model_fd;

	char model_basename_tmp[256];
	uint32_t drawable_addr;
	int cur_drawable = 0;

	if (!strcmp(extension, "xft")) {
		drawable_addr = _ADDR(get_i32_big(&xdr_buf[0x30]));
		model_fd = fopen(model_name, "wb");
		dump_drawable(model_fd, xdr_buf, drawable_addr, model_basename);
		fclose(model_fd);
		printf("Wrote %s\n", model_name);
	} else if (!strcmp(extension, "xdd")) {
		while ((drawable_addr = xdd_get_next_drawable(xdr_buf)) != 0) {
			sprintf(model_basename_tmp, "%s_%i", model_basename, cur_drawable++);
			sprintf(model_name, "%s.obj", model_basename_tmp);
			model_fd = fopen(model_name, "wb");
			dump_drawable(model_fd, xdr_buf, drawable_addr, model_basename_tmp);
			fclose(model_fd);
			printf("Wrote %s\n", model_name);
		}
	} else if (!strcmp(extension, "xdr")) {
		model_fd = fopen(model_name, "wb");
		dump_drawable(model_fd, xdr_buf, 0x10, model_basename);
		fclose(model_fd);
		printf("Wrote %s\n", model_name);
	} else {
		printf("unrecognized extension %s\n", extension);
		return 1;
	}
	return 0;
}

/* ptr is at 0x28, count is at 0x2C */
uint32_t xdd_get_next_drawable(char* xdr_buf) {
	static int next_drawable = 0;

	uint32_t drawable_tbl = _ADDR(get_i32_big(&xdr_buf[0x28]));
	int num_drawables = get_i16_big(&xdr_buf[0x2C]);
	if (next_drawable >= num_drawables) {
		return 0; // no more drawables;
	}
    return _ADDR(get_i32_big(&xdr_buf[drawable_tbl + ((next_drawable++) * 4)]));
}

void dump_drawable(FILE* model_fd, char* xdr_buf, uint32_t drawable_addr, char* model_basename) {
	uint32_t model_addr = _ADDR(get_i32_big(&xdr_buf[drawable_addr + 0x40]));
	uint32_t model_tbl_ptr = _ADDR(get_i32_big(&xdr_buf[model_addr]));
	uint16_t model_count = get_i16_big(&xdr_buf[model_addr + 4]);

	uint32_t idx_ofs = 1;

	/* parse models */
	for (int i = 0; i < model_count; i++) {
		uint32_t model_ptr = _ADDR(get_i32_big(&xdr_buf[model_tbl_ptr+(i*4)]));
		uint32_t mesh_tbl_ptr = _ADDR(get_i32_big(&xdr_buf[model_ptr+(1*4)]));
		uint16_t mesh_count = get_i16_big(&xdr_buf[model_ptr+(2*4)]);

		fprintf(model_fd, "o %s_%i\n", model_basename, i);

		/* parse meshes */
		for (int j = 0; j < mesh_count; j++) {
			uint32_t mesh_ptr = _ADDR(get_i32_big(&xdr_buf[mesh_tbl_ptr+(j*4)]));

			uint32_t vbuf_ptr = _ADDR(get_i32_big(&xdr_buf[mesh_ptr+(3*4)]));
			uint32_t ibuf_ptr = _ADDR(get_i32_big(&xdr_buf[mesh_ptr+(7*4)]));

			uint16_t vbuf_stride = get_i16_big(&xdr_buf[vbuf_ptr+(1*4)]);
			uint16_t ibuf_stride = 2*3;

			uint32_t idx_count = get_i32_big(&xdr_buf[mesh_ptr+(11*4)]);
			uint32_t tri_count = get_i32_big(&xdr_buf[mesh_ptr+(12*4)]);
			uint16_t vert_count = get_i16_big(&xdr_buf[mesh_ptr+(13*4)]);

			uint32_t vbuf_data_ptr = G_ADDR(get_i32_big(&xdr_buf[vbuf_ptr] + 8));
			uint32_t ibuf_data_ptr = G_ADDR(get_i32_big(&xdr_buf[ibuf_ptr] + 8));

			fprintf(model_fd, "g %s_%i_%i\n", model_basename, i, j);

			/* parse vertex buffer */
			for (int k = 0; k < vert_count; k++) {
				float x, y, z, w, u, v;
				x = y = z = w = 0.0f;
				x = get_f32(&xdr_buf[vbuf_data_ptr+(vbuf_stride*k)+(0*4)]);
				y = get_f32(&xdr_buf[vbuf_data_ptr+(vbuf_stride*k)+(1*4)]);
				z = get_f32(&xdr_buf[vbuf_data_ptr+(vbuf_stride*k)+(2*4)]);
				u = get_f16(&xdr_buf[vbuf_data_ptr+(vbuf_stride*k)+(5*4)]);
				v = get_f16(&xdr_buf[vbuf_data_ptr+(vbuf_stride*k)+(5*4)+2]);
				v = (-v)+1;
				fprintf(model_fd, "v %f %f %f\n", x, y, z);
				fprintf(model_fd, "vt %f %f\n", u, v);
			}

			/* parse index buffer */
			for (int k = 0; k < tri_count; k++) {
				uint16_t p0, p1, p2;
				p0 = p1 = p2 = 0;
				p0 = get_i16_big(&xdr_buf[ibuf_data_ptr+(ibuf_stride*k)+(0*2)]) + idx_ofs;
				p1 = get_i16_big(&xdr_buf[ibuf_data_ptr+(ibuf_stride*k)+(1*2)]) + idx_ofs;
				p2 = get_i16_big(&xdr_buf[ibuf_data_ptr+(ibuf_stride*k)+(2*2)]) + idx_ofs;
				fprintf(model_fd, "f %i/%i %i/%i %i/%i\n", p0, p0, p1, p1, p2, p2);
			}
			idx_ofs += vert_count;
		}
	}
}
