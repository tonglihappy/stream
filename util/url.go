package util

import (
	//	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

type UrlSignInfo struct {
	TimeExpiredSecond int
	Md5Err            bool
}

func JoinUrlsArgs(urls []string, args map[string]string) []string {
	for i := 0; i < len(urls); i++ {
		for k, v := range args {
			sep := "?"
			if strings.Contains(urls[i], "?") {
				sep = "&"
			}

			urls[i] = urls[i] + sep + k + "=" + v
		}
	}

	return urls
}

func GetUrl(name string, host string, unique_name string, addr string, usi *UrlSignInfo) ([]string, error) {
	var si StreamUrlInfo

	si.Stream = name
	si.UniqueName = unique_name
	si.Host = host
	si.ServerIp = addr

	return StreamUrl(&si)
}

func GetAllEdgesUrls(name string, host string, unique_name string, time_err bool, md5_err bool) ([]string, error) {
	var urls []string
	var urls_add []string

	var si StreamUrlInfo

	si.Stream = name
	si.UniqueName = unique_name
	si.Host = host
	si.ServerIp = Edge1Addr
	si.ErrorTime = time_err
	si.ErrorMd5 = md5_err

	urls, err := StreamUrl(&si)
	if err != nil {
		return nil, err
	}

	//urls = append(urls, url)

	si.ServerIp = Edge2Addr
	urls_add, err = StreamUrl(&si)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(urls_add); i++ {
		urls = append(urls, urls_add[i])
	}
	//urls = append(urls, url)
	/*
		si.ServerIp = Edge3Addr
		url, err = StreamUrl(&si)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)

		si.ServerIp = Edge3Addr
		url, err = StreamUrl(&si)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	*/
	return urls, nil
}

func GetAllEdgePushUrls(name string, host string, unique_name string, time_err bool, md5_err bool) ([]string, error) {
	return GetAllEdgesUrls(name, host, unique_name, time_err, md5_err)
}

func GetRtmpUrl(addr string, host string, name string) string {
	return "rtmp://" + addr + "/" + host + "/live/" + name
}

func GetHdlUrl(addr string, host string, name string) string {
	return "http://" + addr + "/" + host + "/live/" + name + ".flv"
}

func GetOnlineStream(unique_name string, app string) string {
	var stream_name string
	url := "http://" + unique_name + DashboardSuffix
	paths := "app." + app

	str := GetHttpBody(url)

	mtok := gjson.Get(str, paths)
	var count int
	mtok.ForEach(func(key, value gjson.Result) bool {
		count++
		if count == 3 {
			stream_name = key.String()
			return true
		}
		return true
	})
	return stream_name
}
