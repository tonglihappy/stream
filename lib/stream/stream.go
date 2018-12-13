package stream

//#cgo LDFLAGS: -L../lib/stream/cgo/lib -lffmpeg
//#cgo LDFLAGS: -L./cgo/lib -lavformat -lavcodec -lswscale -lavutil -lavfilter -lswresample -lavdevice -lpostproc -lx264 -ldl -lm -lrt -lpthread -lstdc++
//#cgo CFLAGS: -I ./cgo/include
//#include "main_publish.h"
//#include "main_pull.h"
import "C"

import (
    "fmt"
    "time"
    "unsafe"
)

func StartPushStream(filename string, url string) unsafe.Pointer {
    return C.start_push_stream(C.CString(filename), C.CString(url))
}

func StopStream(stream unsafe.Pointer) string {
    ret := C.stop_stream(unsafe.Pointer(stream))
    return C.GoString(ret)
}

func GetStreamUrl(stream unsafe.Pointer) string {
    ret := C.get_stream_url(unsafe.Pointer(stream))
    return C.GoString(ret)
}

func GetStreamErr(stream unsafe.Pointer) string {
    ret := C.get_stream_err(unsafe.Pointer(stream))
    return C.GoString(ret)
}

func PushStream(filename string, url string) int {
    cf := C.CString(filename)
    curl := C.CString(url)

    ret := int(C.push_stream(cf, curl))
    //fmt.Println(i)
    return ret
}

func PullStream(url string, infi int) int {
    return int(C.pull_stream(C.CString(url), C.int(infi)))
}

func KillStream() {
    C.kill_all_stream()
}

func ThreadKill() {
    //    C.thread_kill()
}

func main() {
    for i := 0; i < 10; i++ {
        go PushStream("/root/livecdn_autotest/src/media/ss.flv", "rtmp://14.17.124.21/videoqa.uplive.ks-cdn.com/live/test011")
        time.Sleep(2 * time.Second)
        ret := PullStream("rtmp://14.17.124.21/videoqa.rtmplive.ks-cdn.com/live/test011", 0)
        if ret == 200 {
            fmt.Println("pull stream success")
        } else {
            fmt.Println("pull stream failed")
        }
        KillStream()
        time.Sleep(1 * time.Second)
    }
}
