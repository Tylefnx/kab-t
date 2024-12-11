#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ncurses.h>
#include <libwebsockets.h>
#include "parson.h"

static int interrupted;
static struct lws *client_wsi;

void show_question(WINDOW *win, const char *question, const int *choices, int num_choices) {
    werase(win);
    box(win, 0, 0);
    mvwprintw(win, 1, 1, "Question: %s", question);

    for (int i = 0; i < num_choices; i++) {
        mvwprintw(win, i + 3, 1, "%d. %d", i + 1, choices[i]);
    }
    wrefresh(win);
}

void handle_message(const char *message) {
    WINDOW *win = newwin(15, 50, 5, 5);
    box(win, 0, 0);
    wrefresh(win);

    JSON_Value *root_value = json_parse_string(message);
    if (json_value_get_type(root_value) != JSONObject) {
        endwin();
        fprintf(stderr, "JSON Parse Error: %s\n", message);
        json_value_free(root_value);
        return;
    }

    JSON_Object *root_object = json_value_get_object(root_value);

    const char *status = json_object_get_string(root_object, "status");
    if (status && strcmp(status, "waiting") == 0) {
        // Handle waiting status
        printf("Status: waiting\n");
    } else if (status && strcmp(status, "quiz") == 0) {
        // Handle quiz question
        JSON_Object *question_object = json_object_get_object(root_object, "question");
        const char *question = json_object_get_string(question_object, "text");
        JSON_Array *choices_array = json_object_get_array(question_object, "choices");

        if (!question || !choices_array) {
            endwin();
            fprintf(stderr, "JSON Format Error: %s\n", message);
            json_value_free(root_value);
            return;
        }

        int num_choices = json_array_get_count(choices_array);
        int choices[num_choices];

        for (int i = 0; i < num_choices; i++) {
            choices[i] = (int)json_array_get_number(choices_array, i);
        }

        show_question(win, question, choices, num_choices);
    } else if (status && strcmp(status, "answer") == 0) {
        // Handle answer status
        int correct_answer = (int)json_object_get_number(root_object, "correct_answer");
        printf("Correct Answer: %d\n", correct_answer);
    } else {
        endwin();
        fprintf(stderr, "JSON Format Error: %s\n", message);
    }

    json_value_free(root_value); // JSON nesnesini serbest bırak
}

static int callback_client(struct lws *wsi, enum lws_callback_reasons reason,
                           void *user, void *in, size_t len)
{
    switch (reason) {
        case LWS_CALLBACK_CLIENT_ESTABLISHED:
            printf("Client connected\n");
            break;
        case LWS_CALLBACK_CLIENT_RECEIVE:
            handle_message((char *)in);
            break;
        case LWS_CALLBACK_CLIENT_WRITEABLE:
            // WebSocket üzerinden yazma işlemleri buraya eklenir.
            break;
        case LWS_CALLBACK_CLIENT_CONNECTION_ERROR:
        case LWS_CALLBACK_CLIENT_CLOSED:
            interrupted = 1;
            break;
        default:
            break;
    }
    return 0;
}

static const struct lws_protocols protocols[] = {
    { "example-protocol", callback_client, 0, 4096, },
    { NULL, NULL, 0, 0 } /* terminator */
};

int main() {
    struct lws_context_creation_info info;
    struct lws_client_connect_info ccinfo;
    struct lws_context *context;

    memset(&info, 0, sizeof(info));
    info.port = CONTEXT_PORT_NO_LISTEN;
    info.protocols = protocols;

    context = lws_create_context(&info);
    if (!context) {
        fprintf(stderr, "lws init failed\n");
        return -1;
    }

    memset(&ccinfo, 0, sizeof(ccinfo));
    ccinfo.context = context;
    ccinfo.address = "localhost";
    ccinfo.port = 8080;
    ccinfo.path = "/ws";
    ccinfo.host = lws_canonical_hostname(context);
    ccinfo.origin = "origin";
    ccinfo.protocol = protocols[0].name;

    client_wsi = lws_client_connect_via_info(&ccinfo);
    if (!client_wsi) {
        fprintf(stderr, "Client connect failed\n");
        lws_context_destroy(context);
        return -1;
    }

    initscr();
    cbreak();
    noecho();

    const char *question = "Waiting for question..."; // Başlangıçta bekleme durumu
    const char *choices[] = {""}; // Başlangıçta boş seçenekler
    int num_choices = 0;

    while (!interrupted) {
        lws_service(context, 1000);

        // Kullanıcıdan giriş almak ve soruları göstermek için burayı güncelleyebilirsiniz.
    }

    endwin();
    lws_context_destroy(context);

    return 0;
}
