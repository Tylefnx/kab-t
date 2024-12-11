#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ncurses.h>
#include "functions.h"

WINDOW *win;
char player_id[51];
int interrupted = 0;
struct lws *client_wsi;
struct lws_protocols protocols[] = {
    {
        "example-protocol",
        callback_client,
        0,
        4096,
        0,
        NULL, NULL, 0
    },
    {NULL, NULL, 0, 0, 0, NULL, NULL, 0} 
};

int main()
{
    struct lws_context *context;

    initscr();
    cbreak();
    noecho();

    // Kullanıcıdan player_id al
    mvprintw(0, 0, "Enter your Player Name: ");
    echo();
    getnstr(player_id, sizeof(player_id) - 1);
    noecho();

    win = newwin(15, 50, 5, 5);
    box(win, 0, 0);
    wrefresh(win);

    // Log seviyesini azalt
    lws_set_log_level(LLL_ERR | LLL_WARN, NULL);

    if (connect_and_join_queue(&context, "localhost", 8080, "/ws", protocols[0].name) != 0)
    {
        endwin();
        return -1;
    }

    while (!interrupted)
    {
        lws_service(context, 1000);
    }

    delwin(win);
    endwin();
    lws_context_destroy(context);

    return 0;
}
