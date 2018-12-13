package util

import (
	"encoding/hex"
	"fmt"
	"github.com/flv.go/flv"
	//"io/ioutil"
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

type FlvCheck struct {
	hc                 *HttpRespCheckInfo
	FirstDataMustVideo bool
	FirstVideoMustKey  bool
	FirstDataMustAudio bool
	NoVideo            bool
	VideoHeaderCnt     int
	NoAudio            bool
	AudioHeaderCnt     int
	VideoCnt           int
	AudioCnt           int
	ReadCnt            int
	DtsStartFrom_0     bool
	KsP2pHint          bool
	KsP2pRange         bool
	MediaSei           bool
	HaveSei            bool
	Url                string
}

type FlvInfo struct {
	totalCnt       int
	videoCnt       int
	audioCnt       int
	metaCnt        int
	videoHeaderCnt int
	audioHeaderCnt int
}

func NewFlvChecker() *FlvCheck {
	checker := &FlvCheck{}

	checker.hc = NewHttpChecker()
	checker.VideoHeaderCnt = 1
	checker.AudioHeaderCnt = 1

	return checker
}

func (check *FlvCheck) AssertRespHeaderCheck(t *testing.T, resp *http.Response) {
	check.hc.AssertRespHeaderCheck(t, resp)
}

func (check *FlvCheck) AssertRespBodyCheck(t *testing.T, resp *http.Response) {
	count := 20
	buf := make([][]byte, 4096)
	videoDtsStartFrom_0 := false
	audioDtsStartfrom_0 := false
	//body := resp.Body

	for i := 0; i < count; i++ {
		body := make([]byte, 2048)
		n, err := resp.Body.Read(body)
		if err != nil && err != io.EOF {
			fmt.Println("get flv stream failed")
		}

		if n == 0 {
			break
		}
		buf[i] = body[:n]
	}
	sep := []byte("")

	body := bytes.NewReader(bytes.Join(buf, sep))
	//fmt.Println(bytes.Join(buf, sep))
	flvReader := flv.NewReader(body)

	info := FlvInfo{}
	var fr flv.Frame
	_, err := flvReader.ReadHeader()
	if err != nil {
		TFatal(t, fmt.Sprintf("read header failed: %s", err))
	}

	for true {
		fr, err = flvReader.ReadFrame()
		if err != nil && err != io.EOF {
			TFatal(t, fmt.Sprintf("read frame failed: %s", err))
		}

		if err == io.EOF {
			break
		}

		info.totalCnt += 1
		if fr.GetType() == flv.TAG_TYPE_VIDEO {
			if fr.IsHeader() {
				info.videoHeaderCnt += 1
			} else {
				if check.MediaSei {
					body := *fr.GetBody()
					encodedStr := hex.EncodeToString(body[5:13])
					//fmt.Println("", encodedStr)
					if check.HaveSei {
						if encodedStr != "000000090600000f" {
							TFatal(t, fmt.Sprintf("no sei"))
						}
					} else {
						if encodedStr == "000000090600000f" {
							TFatal(t, fmt.Sprintf("have sei"))
						}
					}
				}

				if check.NoVideo {
					TFatal(t, fmt.Sprintf("no video but got video"))
				}

				if check.DtsStartFrom_0 {
					if fr.GetDts() == 0 && !videoDtsStartFrom_0 {
						videoDtsStartFrom_0 = true
					}
				}

				if info.audioCnt == 0 {
					if check.FirstDataMustAudio {
						TFatal(t, fmt.Sprintf("first data must audio got video!"))
					}
				}

				info.videoCnt += 1

				if info.videoCnt == 1 {
					if !fr.IsKeyFrame() {
						TFatal(t, fmt.Sprintf("first video must key frame!"))
					}
				}
			}
		} else if fr.GetType() == flv.TAG_TYPE_AUDIO {
			if fr.IsHeader() {
				info.audioHeaderCnt += 1
			} else {
				if check.NoAudio {
					TFatal(t, fmt.Sprintf("no audio but got audio"))
				}

				if check.DtsStartFrom_0 {
					if fr.GetDts() == 0 && !audioDtsStartfrom_0 {
						audioDtsStartfrom_0 = true
					}
				}

				if info.videoCnt == 0 {
					if check.FirstDataMustVideo {
						TFatal(t, fmt.Sprintf("first data must video but got audio!"))
					}
				}
				info.audioCnt += 1
			}

		} else if fr.GetType() == flv.TAG_TYPE_META {
			info.metaCnt += 2
			//body := *fr.GetBody()
			//fmt.Println("meta")
			//fmt.Println(fr.String())
			if check.KsP2pRange {
				KsP2pStreamParser(t, fr.String(), check.Url)
			}
		}

		if info.totalCnt >= check.ReadCnt {
			break
		}
	}
	if check.DtsStartFrom_0 {
		if !audioDtsStartfrom_0 || !videoDtsStartFrom_0 {
			TFatal(t, fmt.Sprintf("audo or video dts not start from 0"))
		}
	}
	if info.videoHeaderCnt != check.VideoHeaderCnt {
		TFatal(t, fmt.Sprintf("except %d video header but got %d video header", check.VideoHeaderCnt, info.videoHeaderCnt))
	}

	if check.KsP2pHint {
		if info.audioHeaderCnt != 1 && info.videoHeaderCnt != 1 {
			TFatal(t, fmt.Sprintf("p2p hint audio header cnt %d video cnt %d", info.audioHeaderCnt, info.videoHeaderCnt))
		}
	}

	if info.audioHeaderCnt != check.AudioHeaderCnt {
		TFatal(t, fmt.Sprintf("except %d audio header but got %d audio header", check.AudioHeaderCnt, info.audioHeaderCnt))
	}
}

func KsP2pStreamParser(t *testing.T, fr string, url string) {
	var push string
	var have []string
	stream := strings.Split(fr, "\t")
	if len(stream) > 0 {
		index_body := strings.Split(stream[4], ";")
		for i := 0; i < len(index_body); i++ {
			if strings.Contains(index_body[i], "stream=") {
				push = strings.Replace(strings.TrimSuffix(index_body[i], "e+12"), ".", "", 1)
				//fmt.Println(push)
			}

			if strings.Contains(index_body[i], "have=") {
				have = strings.Split(index_body[i], ",")
				for i := 1; i < len(have); i++ {
					have[i] = "have=" + have[i]
				}
			}
			if len(push) > 0 && len(have) > 0 {
				DoRangeRequest(t, push, have, url)
			}
		}
	} else {
		TFatal(t, fmt.Sprintf("p2p hint request return error %s", fr))
	}
}

func DoRangeRequest(t *testing.T, push string, have []string, url_str string) {
	client := &http.Client{}
	for i := 0; i < len(have); i++ {
		index := strings.TrimPrefix(strings.Split(have[i], "/")[0], "have=")
		begin := strings.Split(strings.Split(have[i], "/")[1], "-")[0]
		end := strings.Split(strings.Split(have[i], "/")[1], "-")[1]
		push := strings.TrimPrefix(push, "stream=")
		url := fmt.Sprintf("%s&begin=%s&end=%s&index=%s&push=%s", url_str, begin, end, index, push)

		resp, err := client.Get(url)
		defer resp.Body.Close()
		if err != nil {
			TFatal(t, fmt.Sprintf("p2p range request return error %s", url))
		}
		for {

		}
	}
}
