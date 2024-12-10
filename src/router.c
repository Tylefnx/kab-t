#include <microhttpd.h>
#include <string.h>
#include <stdio.h>
#include "question.h"
#include "leaderboard.h"

static int handle_request(void *cls, struct MHD_Connection *connection,
                          const char *url, const char *method,
                          const char *version, const char *upload_data,
                          size_t *upload_data_size, void **con_cls) {
    const char *response = NULL;
    int ret;

    if (strcmp(url, "/generate") == 0 && strcmp(method, "GET") == 0) {
        response = generate_questions();
    } else if (strcmp(url, "/leaderboard") == 0 && strcmp(method, "GET") == 0) {
        response = get_leaderboard();
    } else {
        response = "{\"error\":\"Endpoint not found\"}";
    }

    struct MHD_Response *res = MHD_create_response_from_buffer(strlen(response), 
                                                               (void *)response, 
                                                               MHD_RESPMEM_MUST_COPY);
    ret = MHD_queue_response(connection, MHD_HTTP_OK, res);
    MHD_destroy_response(res);
    return ret;
}

void start_server(int port) {
    struct MHD_Daemon *daemon = MHD_start_daemon(MHD_USE_SELECT_INTERNALLY, port,
                                                 NULL, NULL, &handle_request,
                                                 NULL, MHD_OPTION_END);
    if (daemon == NULL) {
        fprintf(stderr, "Failed to start server\n");
        exit(1);
    }

    printf("Server running on port %d...\n", port);
    getchar();
    MHD_stop_daemon(daemon);
}
