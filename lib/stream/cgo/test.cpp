
#include <stdio.h>
#include <unistd.h>
#include "main_publish.h"

int main(int argc, char *argv[]) {
    printf("start test!\n");

    char  url[256];
    void *ps[100];

    for (int i = 0; i < 100; i++) {
        sprintf(url, "rtmp://112.84.131.160/livecdn-qa-online.uplive.kscvbu.cn/live/ffmpeg_test_%d", i);
        ps[i] = start_push_stream("/root/livecdn_autotest/src/media/ss.flv", url);
    }

    sleep(30);

    for (int i = 0; i < 100; i++) {
        if (ps[i] == NULL) {
            continue;
        }

        printf("finish: %d, %s\n", i, stop_stream(ps[i]));
    }

    /*
    void *ts = start_push_stream("/root/livecdn_autotest/src/media/ss.flv", "rtmp://112.84.131.160/livecdn-qa-online.uplive.kscvbu.cn/live/ffmpeg_test");
    void *ts2 = start_push_stream("/root/livecdn_autotest/src/media/ss.flv", "rtmp://112.84.131.160/livecdn-qa-online.uplive.kscvbu.cn/live/ffmpeg_test2");

    if (ts == NULL || ts2 == NULL) {
        printf("push stream error!\n");
        return -1;
    }

    sleep(10);

    printf("finish ts %s\n", stop_stream(ts));
    printf("finish ts2 %s\n", stop_stream(ts2));

    ts = start_push_stream("/root/livecdn_autotest/src/media/ss.flv", "rtmp://112.84.131.160/livecdn-qa-online.uplive.kscvbu.cn/live/ffmpeg_test");
    sleep(10);

    if (ts == NULL) {
        printf("push stream error!\n");
    }
    
    printf("finish ts3 %s\n", stop_stream(ts));
    */

    return 0;
}
