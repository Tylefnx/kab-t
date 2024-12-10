#include <stdio.h>
#include <stdlib.h>
#include "router.h"

int main() {
    printf("Kahoot Backend is starting...\n");
    start_server(8080);
    return 0;
}
