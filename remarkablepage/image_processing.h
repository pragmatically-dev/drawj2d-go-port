#ifndef IMAGE_PROCESSING_H
#define IMAGE_PROCESSING_H

#include <stdlib.h>
#include <stdbool.h>


typedef struct {
    float *lines;
    int size;
} LineList;

LineList handle_new_file(const char *directory, const char *filename);

#endif
