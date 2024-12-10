#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "question.h"

char *generate_questions() {
    static char response[512];
    snprintf(response, sizeof(response),
             "{\"questions\":[{\"q\":\"%d + %d\", \"id\":1}, {\"q\":\"%d + %d\", \"id\":2}]}",
             rand() % 1000, rand() % 1000, rand() % 1000, rand() % 1000);
    return response;
}
