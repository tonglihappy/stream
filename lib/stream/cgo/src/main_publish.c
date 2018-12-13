/*================================================================
*   Copyright (C) 2017 KingSoft Cloud Ltd. All rights reserved.
*   
*   文件名称：main_publish.c
*   创 建 者：tongli
*   创建日期：2017年10月09日
*   描    述：
*
================================================================*/
#include "main_publish.h"
#define __STDC_CONSTANT_MACROS
#include "stream.h"
#include <pthread.h>
#include <unistd.h>
#include <iostream>

int pull_stream(const char* url){

}

int push_stream(char* filename, char* url){
	int ret;
	init();
	if(ret = create_stream(filename, url) < 0)
		std::cout << "create stream failed" << std::endl;

	deinit();
	return ret;
}

void kill_all_stream(){
	close_stream();	
}


void *push_stream_thread(void *arg) {
    Stream *stream = (Stream *)arg;
    int ret = 0;

	if(ret = stream->open_input_file(stream->m_file) < 0){
        if (stream->m_err == "") {
            stream->m_err = std::string("open input file failed: ") + stream->m_file;
        }
        return NULL;
	}

	if(ret = stream->open_output_file(stream->m_url) < 0){
        if (stream->m_err == "") { 
            stream->m_err = std::string("open output file failed") + stream->m_url;
        }
        return NULL;
	}

	stream->transcode();

    return NULL;
}


void *start_push_stream(char *filename, char *url) {
    init();
    Stream *stream = new Stream(filename, url);

    pthread_t   td;
    int ret = 0;

    ret = pthread_create(&td, NULL, &push_stream_thread, stream);
    if (ret != 0) {
        stream->m_err = "create thread failed!";
        goto error;
    }

    while (true) {
        if (stream->m_err != "") {
            break;
        }

        if (stream->get_written_pkt_cnt() > 200) {
            break;
        }

        if (stream->m_loop_stop) {
            if (stream->m_err == "") {
                stream->m_err = "stream terminal premature";
            }
            break;
        }

        sleep(1);
    }

    if (stream->m_err != "") {
        goto error;
    }

    return stream;


error:
	return stream; 
}

const char *get_stream_url(void *opaque) {
    Stream *stream = (Stream *)opaque;

    return stream->m_url;
}

const char *get_stream_err(void *opaque) {
    Stream *stream = (Stream *)opaque;

    return stream->m_err.c_str();
}

const char *stop_stream(void *opaque) {
    Stream *stream = (Stream *)opaque;

    std::string     err("");

    if (!stream->m_loop_start || stream->m_err != "") {
        err = get_stream_err(opaque);
        delete stream;

        return err.c_str();
    } 

    stream->m_exit = true;

    while (!stream->m_loop_stop) {
        sleep(1);
    }

    err = get_stream_err(opaque);
    delete stream;

    return err.c_str();
}


