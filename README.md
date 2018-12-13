
#简介：<br>
本工程是流媒体自动化测试的基本框架。包含基本推拉流，http/rtmp/websocket等协议的校验，flv的基础校验。可用作流媒体服务器自动化测试的基础框架。
代码作为一些参考，因为涉及到公司的一些不能透漏的信息，只把底层的基础拿出来分享。<br>

#举例子：
##1.测试基本推拉流
`
func TestVideoqaBaseStream(t *testing.T) {
        name := "cdn_base_test1"
        defer util.AssertStopStream(t, util.AssertPushStream(t, util.MediaFile, getPushUrl(name, util.Edge1Addr)))
        util.AssertMultiNormalHdl(t, getUrls(name, false, false), 10240, true)
}
`

2.测试时间戳从0开始
func TestVideoqaDtsStartFrom_0(t *testing.T) {
    name := "cdn_base_test2"

    defer util.AssertStopStream(t, util.AssertPushStream(t, util.MediaFileDtsNotStartFrom_0, getPushUrl(name, util.Edge1Addr)))
    treq := util.NewHttpReq()
    treq.Timeout = 20 * time.Second
    treq.Url = "http://domain/app/streamname.flv"

    flvChecker := util.NewFlvChecker()
    flvChecker.DtsStartFrom_0 = true
    flvChecker.ReadCnt = 100

    util.AssertHttpRequest(t, treq, flvChecker)
}

3.http协议在拉不到流的状况下返回给client的状态码的校验

func TestVideoqaCheckPullStreamFailedStatusCode(t *testing.T) {
    name := "test_snipper"
    treq := util.NewHttpReq()
    treq.Timeout = 30 * time.Second
    treq.Url = "http://domain/app/streamname.flv"

    httpChecker := util.NewHttpChecker()
    httpChecker.StatusCode = 404

    util.AssertHttpRequest(t, treq, httpChecker)
}
