package web

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	pb "github.com/cheggaaa/pb/v3"
)

var client = &http.Client{}

func Download(url string, target string) bool {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// 添加常见的浏览器 User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false
	}

	if response.StatusCode != 200 {
		fmt.Println("Error status while downloading", url, "-", response.StatusCode)
		return false
	}
	defer response.Body.Close()

	output, err := os.Create(target)
	if err != nil {
		fmt.Println("Error while creating", target, "-", err)
		return false
	}
	defer output.Close()

	// 创建进度条 - 使用默认模板并设置总字节数
	bar := pb.New(int(response.ContentLength)).
		Set(pb.Bytes, true).                   // 显示字节单位
		SetRefreshRate(time.Millisecond * 10). // 刷新频率
		SetWidth(80)                           // 宽度

	// 如果需要更多自定义，可以使用完整模板
	// bar := pb.Full.New(int(response.ContentLength)).
	//     Set(pb.Bytes, true).
	//     SetRefreshRate(time.Millisecond * 10).
	//     SetWidth(80)

	bar.Start()

	proxyWriter := bar.NewProxyWriter(output)

	_, err = io.Copy(proxyWriter, response.Body)

	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false
	}
	bar.Finish()

	return true
}

func GetRemoteTextFile(url string) (string, error) {
	response, httperr := client.Get(url)
	if httperr != nil {
		return "", errors.New(fmt.Sprintf("\nCould not retrieve %s.\n\n%s\n", url, httperr.Error()))
	} else {
		defer response.Body.Close()
		contents, readerr := ioutil.ReadAll(response.Body)
		if readerr != nil {
			return "", errors.New(fmt.Sprintf("%s", readerr))
		}
		return string(contents), nil
	}
}

func GetPython(download string, v string, url string) (string, bool) {
	fileName := filepath.Join(download, fmt.Sprintf("%s.zip", v))
	os.Remove(fileName)
	if url == "" {
		//No url should mean this version/arch isn't available
		fmt.Printf("Python %s isn't available right now.", v)
	} else {
		fmt.Printf("Downloading python version %s...\n", v)
		if Download(url, fileName) {
			fmt.Println("Complete")
			return fileName, true
		} else {
			return "", false
		}
	}
	return "", false
}
