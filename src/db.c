#include <sqlite3.h>
#include <stdio.h>
#include <stdlib.h>
#include "db.h"

sqlite3 *db;

void init_db() {
    if (sqlite3_open("kahoot.db", &db) != SQLITE_OK) {
        fprintf(stderr, "Cannot open database: %s\n", sqlite3_errmsg(db));
        exit(1);
    }
    const char *sql = "CREATE TABLE IF NOT EXISTS leaderboard (id INTEGER PRIMARY KEY, name TEXT, score INTEGER);";
    char *err_msg = NULL;
    if (sqlite3_exec(db, sql, 0, 0, &err_msg) != SQLITE_OK) {
        fprintf(stderr, "SQL error: %s\n", err_msg);
        sqlite3_free(err_msg);
        sqlite3_close(db);
        exit(1);
    }
}

void close_db() {
    sqlite3_close(db);
}
