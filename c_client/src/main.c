#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ncurses.h>
#include <libwebsockets.h>
#include <curl/curl.h>
#include "parson.h"
#include "functions.h" // Yeni eklediğiniz .h dosyası

struct lws *client_wsi;

static const struct lws_protocols protocols[] = {
    {
        "example-protocol",
        callback_client, // functions.c dosyasına taşıdığımız fonksiyon
        0,
        4096,
        0, // `id` alanı eklendi
        NULL, NULL, 0
    },
    {NULL, NULL, 0, 0, 0, NULL, NULL, 0} /* terminator */
};

int main()
{
    struct lws_context_creation_info info;
    struct lws_client_connect_info ccinfo;
    struct lws_context *context;

    memset(&info, 0, sizeof(info));
    info.port = CONTEXT_PORT_NO_LISTEN;
    info.protocols = protocols;

    context = lws_create_context(&info);
    if (!context)
    {
        fprintf(stderr, "lws init failed\n");
        return -1;
    }

    memset(&ccinfo, 0, sizeof(ccinfo));
    ccinfo.context = context;
    ccinfo.address = "localhost";
    ccinfo.port = 8080;
    ccinfo.path = "/ws";
    ccinfo.origin = "origin";
    ccinfo.protocol = protocols[0].name;
    client_wsi = lws_client_connect_via_info(&ccinfo);
    if (!client_wsi)
    {
        fprintf(stderr, "Client connect failed\n");
        lws_context_destroy(context);
        return -1;
    }

    initscr();
    cbreak();
    noecho();

    // Kullanıcıdan player_id al
    mvprintw(0, 0, "Enter your Player ID: ");
    echo();
    getnstr(player_id, sizeof(player_id) - 1); // Maksimum 50 karakter al
    noecho();

    win = newwin(15, 50, 5, 5);
    box(win, 0, 0);
    wrefresh(win);

    while (!interrupted)
    {
        lws_service(context, 1000);
    }

    delwin(win);
    endwin();
    lws_context_destroy(context);

    return 0;
}
