package util

import (
	"fmt"
	rtmp "github.com/tongli/gortmp"
	"testing"
	"time"
)

type RtmpChecker struct {
}

type RtmpReq struct {
	Url        string
	StreamName string
	Savetofile bool
}

type TestOutboundConnHandler struct {
	streamName string
	t          *testing.T
	tcheck     *FlvCheck
}

func (handler *TestOutboundConnHandler) OnStatus(conn rtmp.OutboundConn) {
	fmt.Printf("@@@@@@@@@@@@@Status\n")
}

func (handler *TestOutboundConnHandler) OnClosed(conn rtmp.Conn) {
	fmt.Printf("@@@@@@@@@@@@@Closed\n")
}

func (handler *TestOutboundConnHandler) OnReceived(conn rtmp.Conn, message *rtmp.Message) {

}

func (handler *TestOutboundConnHandler) OnReceivedRtmpCommand(conn rtmp.Conn, command *rtmp.Command) {
	fmt.Printf("ReceviedCommand: %+v\n", command)
}

func (handler *TestOutboundConnHandler) OnStreamCreated(conn rtmp.OutboundConn, stream rtmp.OutboundStream) {
	t := handler.t
	fmt.Printf("Stream created: %d\n", stream.ID())
	err := stream.Play(handler.streamName, nil, nil, nil)
	if err != nil {
		TFatal(t, fmt.Sprintf("Play error: %s", err))
	}
}

func AssertRtmpRequest(t *testing.T, rtmpreq *RtmpReq, tcheck *FlvCheck) {
	testHandler := &TestOutboundConnHandler{}
	testHandler.streamName = rtmpreq.StreamName
	testHandler.t = t
	testHandler.tcheck = tcheck

	obConn, err := rtmp.Dial(rtmpreq.Url, testHandler, 100)
	if err != nil {
		TFatal(t, fmt.Sprintf("Dial error %s", err))
	}

	defer obConn.Close()

	err = obConn.Connect()
	if err != nil {
		TFatal(t, fmt.Sprintf("Connect error %s", err))
	}

	for {
		select {
		case <-time.After(1 * time.Second):
			continue
		}
	}
}
