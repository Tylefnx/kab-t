#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ncurses.h>
#include <time.h>
#include <curl/curl.h>
#include "parson.h"
#include "functions.h"
// Global değişkenlerin tanımlanması
WINDOW *win;
char player_id[51];
int interrupted = 0;

void send_answer(const char *player_id, int question_id, int choice)
{
    CURL *curl;
    CURLcode res;

    curl = curl_easy_init();
    if (curl)
    {
        char postfields[150];
        snprintf(postfields, sizeof(postfields), "{\"player_id\": \"%s\", \"question_id\": %d, \"answer\": %d}", player_id, question_id, choice);

        struct curl_slist *headers = NULL;
        headers = curl_slist_append(headers, "Content-Type: application/json");

        curl_easy_setopt(curl, CURLOPT_URL, "http://localhost:8080/answer");
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, postfields);
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

        res = curl_easy_perform(curl);
        if (res != CURLE_OK)
            fprintf(stderr, "curl_easy_perform() failed: %s\n", curl_easy_strerror(res));

        curl_slist_free_all(headers); // Hafıza sızıntısını önlemek için header'ı serbest bırak
        curl_easy_cleanup(curl);
    }
}


void show_question(const char *question, const int *choices, int num_choices)
{
    werase(win);
    box(win, 0, 0);
    mvwprintw(win, 1, 1, "Question: %s", question);

    for (int i = 0; i < num_choices; i++)
    {
        mvwprintw(win, i + 3, 1, "%d. %d", i + 1, choices[i]);
    }
    mvwprintw(win, num_choices + 4, 1, "Enter the number of your choice and press Enter:");
    wrefresh(win);
}

#include <time.h>

void handle_message(const char *message)
{
    JSON_Value *root_value = json_parse_string(message);
    if (json_value_get_type(root_value) != JSONObject)
    {
        fprintf(stderr, "JSON Parse Error: %s\n", message);
        json_value_free(root_value);
        return;
    }

    JSON_Object *root_object = json_value_get_object(root_value);

    const char *status = json_object_get_string(root_object, "status");
    if (status && strcmp(status, "waiting") == 0)
    {
        werase(win);
        mvwprintw(win, 1, 1, "Waiting for more players...");
        wrefresh(win);
    }
    else if (status && strcmp(status, "quiz") == 0)
    {
        JSON_Object *question_object = json_object_get_object(root_object, "question");
        const char *question = json_object_get_string(question_object, "text");
        JSON_Array *choices_array = json_object_get_array(question_object, "choices");
        int question_id = (int)json_object_get_number(question_object, "id");
        int timeout = 10; // Her sorunun süresi 10 saniye
        int correct_answer = (int)json_object_get_number(question_object, "answer");

        if (!question || !choices_array)
        {
            fprintf(stderr, "JSON Format Error: %s\n", message);
            json_value_free(root_value);
            return;
        }

        int num_choices = json_array_get_count(choices_array);
        int *choices = malloc(num_choices * sizeof(int));
        if (!choices)
        {
            fprintf(stderr, "Memory Allocation Error\n");
            json_value_free(root_value);
            return;
        }

        for (int i = 0; i < num_choices; i++)
        {
            choices[i] = (int)json_array_get_number(choices_array, i);
        }

        show_question(question, choices, num_choices);

        mvwprintw(win, num_choices + 5, 1, "Time remaining: %d seconds", timeout);
        mvwprintw(win, num_choices + 6, 1, "Enter the number of your choice and press Enter:");
        wrefresh(win);

        struct timespec sleep_time = {1, 0}; // 1 saniye
        int answered = 0; // Kullanıcı cevap verdi mi kontrolü

        for (int remaining_time = timeout - 1; remaining_time >= 0; remaining_time--)
        {
            mvwprintw(win, num_choices + 5, 1, "Time remaining: %d seconds", remaining_time);
            wrefresh(win);
            nanosleep(&sleep_time, NULL); // 1 saniye beklet

            nodelay(win, TRUE); // Non-blocking mode
            int ch = wgetch(win);
            nodelay(win, FALSE); // Reset to blocking mode

            if (ch != ERR) // There is a key press
            {
                int choice_index = ch - '1';

                if (choice_index >= 0 && choice_index < num_choices)
                {
                    mvwprintw(win, num_choices + 7, 1, "You selected choice %d: %d", choice_index + 1, choices[choice_index]);
                    send_answer(player_id, question_id, choices[choice_index]);
                    wrefresh(win);
                    answered = 1; // Kullanıcı cevap verdi
                    break;
                }
            }
        }

        if (!answered) // Kullanıcı cevap vermediyse süre dolduğunda doğru cevabı göster
        {
            mvwprintw(win, num_choices + 7, 1, "Time's up! Correct Answer: %d", correct_answer);
        }

        wrefresh(win);
        free(choices);
    }
    else if (status && strcmp(status, "answer") == 0)
    {
        int correct_answer = (int)json_object_get_number(root_object, "correct_answer");
        mvwprintw(win, 15, 1, "Correct Answer: %d", correct_answer);
        wrefresh(win);
    }
    else if (status && strcmp(status, "leaderboard") == 0)
    {
        JSON_Array *scores_array = json_object_get_array(root_object, "scores");
        int num_scores = json_array_get_count(scores_array);

        werase(win);
        mvwprintw(win, 1, 1, "Leaderboard:");

        for (int i = 0; i < num_scores; i++)
        {
            JSON_Object *score_object = json_array_get_object(scores_array, i);
            const char *player_id = json_object_get_string(score_object, "id");
            int score = (int)json_object_get_number(score_object, "score");

            mvwprintw(win, i + 3, 1, "%d. %s: %d", i + 1, player_id, score);
        }
        wrefresh(win);
    }
    else
    {
        fprintf(stderr, "JSON Format Error: %s\n", message);
    }

    json_value_free(root_value); // JSON nesnesini serbest bırak
}


int callback_client(struct lws *wsi, enum lws_callback_reasons reason,
                    void *user, void *in, size_t len)
{
    switch (reason)
    {
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
