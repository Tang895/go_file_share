package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"

	qrcode "github.com/skip2/go-qrcode"
)

var filename string

func sendfile(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "hello world!")
	header := w.Header()
	header.Add("Content-Type:", "application/octet-stream")
	header.Add("Content-Disposition", "attachment;filename="+filename[2:]) //[2:]去掉前面的./
	file, _ := ioutil.ReadFile("./" + filename)
	si, _ := os.Stat("./" + filename)

	leng := strconv.Itoa(int(si.Size()))
	header.Add("Content-Length", leng)
	w.Write(file)
}
func main() {
	//sys_Info := runtime.GOOS //获取计算机系统

	interfaceaddrs, _ := net.InterfaceAddrs()
	ip := ""
	for i := range interfaceaddrs {
		ip = interfaceaddrs[i].String()
		if ip[:8] == "192.168." {
			ip = ip[:len(ip)-3] //去掉掩码
			break
		}
	}
	filename = os.Args[1]

	//加md5
	filename_Byte := md5.Sum([]byte(filename))
	filename_MD5 := fmt.Sprintf("%x", filename_Byte)

	http.HandleFunc("/"+filename_MD5, sendfile)

	dl_URL := "http://" + ip + ":8000/" + filename_MD5

	q, _ := qrcode.New(dl_URL, qrcode.Low)
	art := q.ToString(false)
	fmt.Println(art) //输出二维码
	fmt.Println("filename:" + filename)
	fmt.Println(dl_URL)
	http.ListenAndServe(":8000", nil)
}
