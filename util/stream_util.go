package util

import (
	"fmt"
	stream "lib/stream"
	"testing"
	"time"
	"unsafe"
)

func AssertPushStream(t *testing.T, file string, url string) unsafe.Pointer {
	//    t.Logf("start push stream, file: %s, url: %s", file, url)
	ps := stream.StartPushStream(file, url)
	err := stream.GetStreamErr(ps)
	if err != "" {
		TFatal(t, fmt.Sprintf("push stream failed, err: %s", err))
		stream.StopStream(ps)

		return nil
	}

	/* wait some time */
	time.Sleep(time.Second * 5)

	return ps
}

func AssertPublishStream(t *testing.T, file string, url string, ok bool) unsafe.Pointer {
	//    t.Logf("start push stream, file: %s, url: %s", file, url)
	ps := stream.StartPushStream(file, url)
	err := stream.GetStreamErr(ps)
	//fmt.Println(url)

	if ok && err != "" {
		TFatal(t, fmt.Sprintf("Failed! need ok, push stream failed, err:%s", err))
		stream.StopStream(ps)

		return nil
	}

	if !ok && err == "" {
		TFatal(t, fmt.Sprintf("Failed! need failed, push stream success."))
		stream.StopStream(ps)

		return nil
	}

	/* wait some time */
	time.Sleep(time.Second * 5)

	return ps
}

func AssertStopStream(t *testing.T, ps unsafe.Pointer) {
	//   t.Logf("stop stream %s", stream.GetStreamUrl(ps))
	err := stream.StopStream(ps)
	if err != "" {
		TFatal(t, fmt.Sprintf("push stream failed!, %s", err))
	}
}
