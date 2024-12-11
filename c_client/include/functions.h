#ifndef FUNCTIONS_H
#define FUNCTIONS_H

#include <ncurses.h>
#include <libwebsockets.h>

extern WINDOW *win;
extern char player_id[51];
extern int interrupted;
extern struct lws *client_wsi;
extern struct lws_protocols protocols[];

void send_answer(const char *player_id, int question_id, int choice);
void show_question(const char *question, const int *choices, int num_choices);
void handle_message(const char *message);
int callback_client(struct lws *wsi, enum lws_callback_reasons reason, void *user, void *in, size_t len);
int connect_and_join_queue(struct lws_context **context, const char *address, int port, const char *path, const char *protocol);

#endif // FUNCTIONS_H
