/*================================================================
*   Copyright (C) 2017 KingSoft Ltd. All rights reserved.
*   
*   文件名称：main_publish.h
*   创 建 者：tongli
*   创建日期：2017年10月09日
*   描    述：
*
================================================================*/


#pragma once
#ifdef __cplusplus
extern "C" {
#endif
int push_stream(char* filename, char* url);
void kill_all_stream();

void *start_push_stream(char *filename, char *url);
const char *stop_stream(void *opaque);
const char *get_stream_url(void *opaque);
const char *get_stream_err(void *opaque);

#ifdef __cplusplus
}
#endif
