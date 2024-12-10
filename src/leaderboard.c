#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "leaderboard.h"
#include "db.h"

char *get_leaderboard() {
    static char response[512];
    snprintf(response, sizeof(response),
             "{\"leaderboard\":[{\"name\":\"Player1\", \"score\":100}, {\"name\":\"Player2\", \"score\":90}]}");
    return response;
}
