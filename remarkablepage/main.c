#include <stdio.h>
#include <stdlib.h>
#include <sys/inotify.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <linux/limits.h>
#include <stdbool.h>
#include <omp.h>


#define STB_IMAGE_IMPLEMENTATION
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb_image.h"
#include "stb_image_write.h"

typedef struct
{
    float x;
    float y;
} Point;

typedef struct
{
    float *lines;
    int size;
} LineList;

LineList handle_new_file(const char *directory, const char *filename);
void apply_gaussian_blur(unsigned char *image, int width, int height);
void apply_laplace_filter(unsigned char *image, unsigned char *output, int width, int height);
bool **build_boolean_matrix(unsigned char *image, int width, int height);
LineList GetHorizontalLines(bool **pointMatrix, int imageWidth, int imageHeight);



LineList handle_new_file(const char *directory, const char *filename)
{
    LineList horizontalLines ;
    if (strncmp(filename, "Screenshot", 10) == 0)
    {
        char filepath[PATH_MAX];
        snprintf(filepath, PATH_MAX, "%s/%s", directory, filename);
       
        int width, height, channels;
        unsigned char *image = stbi_load(filepath, &width, &height, &channels, STBI_grey); // Load as grayscale (1 channel)
        if (!image)
        {
            fprintf(stderr, "Error loading image %s\n", filepath);
            return horizontalLines;
        }

        unsigned char *output = (unsigned char *)malloc(width * height);
        if (!output)
        {
            fprintf(stderr, "Error allocating memory for output image\n");
            stbi_image_free(image);
            return horizontalLines;
        }

        apply_gaussian_blur(image, width, height);
        apply_laplace_filter(image, output, width, height);

        bool **bool_matrix = build_boolean_matrix(output, width, height);
        if (!bool_matrix)
        {
            fprintf(stderr, "Error creating boolean matrix\n");
            free(output);
            stbi_image_free(image);
            return horizontalLines;
        }


        // Obtener las líneas horizontales
         horizontalLines = GetHorizontalLines(bool_matrix, width, height);

      /*   // Imprimir las líneas horizontales
        for (int i = 0; i < horizontalLines.size; ++i)
        {
            printf("Line %d: (%.2f, %.2f) to (%.2f, %.2f)\n",
                   i,
                   horizontalLines.lines[i * 4],
                   horizontalLines.lines[i * 4 + 1],
                   horizontalLines.lines[i * 4 + 2],
                   horizontalLines.lines[i * 4 + 3]);
        } */

        // Clean up
        for (int i = 0; i < width; ++i)
        {
            free(bool_matrix[i]);
        }
        free(bool_matrix);
        free(output);
        //free(horizontalLines.lines);
        stbi_image_free(image);
        return horizontalLines;
    }
}

void apply_gaussian_blur(unsigned char *image, int width, int height)
{
    const float a = 1.0 / 16.0;
    const float b = 2.0 / 16.0;
    const float c = 4.0 / 16.0;

    const float kernel[3][3] = {
        {a, b, a},
        {b, c, b},
        {a, b, a}
    };

    unsigned char *temp = (unsigned char *)malloc(width * height);
    if (!temp)
    {
        fprintf(stderr, "Error allocating memory for Gaussian blur\n");
        return;
    }

#pragma omp parallel for collapse(2)
    for (int y = 1; y < height - 1; ++y)
    {
        for (int x = 1; x < width - 1; ++x)
        {
            float sum = 0.0;
            for (int ky = -1; ky <= 1; ++ky)
            {
                for (int kx = -1; kx <= 1; ++kx)
                {
                    int pixel = image[(y + ky) * width + (x + kx)];
                    sum += pixel * kernel[ky + 1][kx + 1];
                }
            }
            temp[y * width + x] = (unsigned char)sum;
        }
    }

    memcpy(image, temp, width * height);
    free(temp);
}

void apply_laplace_filter(unsigned char *image, unsigned char *output, int width, int height)
{
    const int kernel[3][3] = {
    	{1, 4, 1},
	{4, -20, 4},
	{1, 4, 1},};

#pragma omp parallel for collapse(2)
    for (int y = 1; y < height - 1; ++y)
    {
        for (int x = 1; x < width - 1; ++x)
        {
            int sum = 0;
            for (int ky = -1; ky <= 1; ++ky)
            {
                for (int kx = -1; kx <= 1; ++kx)
                {
                    int pixel = image[(y + ky) * width + (x + kx)];
                    sum += pixel * kernel[ky + 1][kx + 1];
                }
            }
            sum = sum > 255 ? 255 : (sum < 0 ? 0 : sum);
            output[y * width + x] = (unsigned char)sum;
        }
    }
}


bool **build_boolean_matrix(unsigned char *image, int width, int height)
{
    bool **bool_img_map = (bool **)malloc(width * sizeof(bool *));
    if (!bool_img_map)
    {
        fprintf(stderr, "Error allocating memory for boolean matrix\n");
        return NULL;
    }

    for (int i = 0; i < width; ++i)
    {
        bool_img_map[i] = (bool *)malloc(height * sizeof(bool));
        if (!bool_img_map[i])
        {
            fprintf(stderr, "Error allocating memory for boolean matrix row\n");
            for (int k = 0; k < i; ++k)
            {
                free(bool_img_map[k]);
            }
            free(bool_img_map);
            return NULL;
        }
    }

    for (int i = 0; i < width; ++i)
    {
        for (int j = 0; j < height; ++j)
        {
            bool_img_map[i][j] = image[j * width + i] > 0;
        }
    }

    return bool_img_map;
}

LineList GetHorizontalLines(bool **pointMatrix, int imageWidth, int imageHeight)
{
    int *list = (int *)malloc(imageWidth * imageHeight * 3 * sizeof(int)); // Maximum possible size
    int listSize = 0;
    int from, to;
    bool isLine, isPixel;

    for (int z = 0; z < imageHeight; z++)
    {
        isLine = false;
        for (int x = 0; x < imageWidth; x++)
        {
            isPixel = pointMatrix[x][z];
            if (isLine)
            {
                if (isPixel)
                {
                    to = x;
                    if (x + 1 == imageWidth)
                    {
                        list[listSize++] = z;
                        list[listSize++] = from;
                        list[listSize++] = to;
                        isLine = false;
                    }
                }
                else
                { // !isPixel, line ended
                    list[listSize++] = z;
                    list[listSize++] = from;
                    list[listSize++] = to;
                    isLine = false;
                }
            }
            else
            { // !isLine, line not started
                if (isPixel)
                {
                    from = x;
                    to = x;
                    if (x + 1 == imageWidth)
                    { // single pixel at last row
                        list[listSize++] = z;
                        list[listSize++] = from;
                        list[listSize++] = to;
                    }
                    else
                    {
                        isLine = true;
                    }
                } // else do nothing
            }
        }
    }

    float *horizontalLines = (float *)malloc(listSize / 3 * 4 * sizeof(float)); // Each line has 4 float values
    int horizontalLinesSize = 0;

    for (int i = 0; i < listSize; i += 3)
    {
        Point startPoint = {(float)list[i + 1], (float)list[i]};

        if (list[i + 1] != list[i + 2])
        { // separate start and end points
            Point endPoint = {(float)list[i + 2], (float)list[i]};

            horizontalLines[horizontalLinesSize++] = startPoint.x;
            horizontalLines[horizontalLinesSize++] = startPoint.y;
            horizontalLines[horizontalLinesSize++] = endPoint.x;
            horizontalLines[horizontalLinesSize++] = endPoint.y;
           //  printf("Line %d: (%.2f, %.2f) to (%.2f, %.2f)\n",
             //       i,
               //     startPoint.x,startPoint.y,endPoint.x,endPoint.y );
        
        }
        else
        { // start and end points are equal, thus it is just a point
            horizontalLines[horizontalLinesSize++] = startPoint.x;
            horizontalLines[horizontalLinesSize++] = startPoint.y;
            horizontalLines[horizontalLinesSize++] = startPoint.x;
            horizontalLines[horizontalLinesSize++] = startPoint.y;
           // printf("POINT %d: (%.2f, %.2f) \n",i,startPoint.x,startPoint.y );
        } 
    }

    free(list);

    LineList result;
    result.lines = horizontalLines;
    result.size = horizontalLinesSize / 4;
    return result;
}