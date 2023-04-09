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
var filesize int
var file_Path string

func sendfile(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "hello world!")
	header := w.Header()
	header.Add("Content-Type", "application/octet-stream")
	header.Add("Content-Disposition", "attachment;filename="+filename)
	file, _ := ioutil.ReadFile(file_Path)

	header.Add("Content-Length", strconv.Itoa(filesize))
	w.Write(file)
	//fmt.Println(w.Header())
}
func main() {
	//sys_Info := runtime.GOOS //获取计算机系统

	interfaceaddrs, _ := net.InterfaceAddrs()
	ip := ""
	var ips []string
	ipItor := 1
	for i := range interfaceaddrs {
		ip = interfaceaddrs[i].String()

		ips = append(ips, ip)
		fmt.Println("[", ipItor, "]", ": ", ip)
		ipItor += 1
		// if ip[:8] == "192.168." {
		// 	ip = ip[:len(ip)-3] //去掉掩码
		// 	break
		// }
	}
	//fmt.Println(ips)
	fmt.Print("Select your interface, 选择网口->")
	var interfaceSelect int
	fmt.Scanln(&interfaceSelect)
	ip = ips[interfaceSelect-1]
	ip = ip[:len(ip)-3] //去掉掩码
	file_Path = os.Args[1]
	file_Info, _ := os.Stat(file_Path)
	filesize = int(file_Info.Size())
	filename = file_Info.Name()
	//加md5
	filename_Byte := md5.Sum([]byte(file_Path))
	filename_MD5 := fmt.Sprintf("%x", filename_Byte)

	http.HandleFunc("/"+filename_MD5, sendfile)

	dl_URL := "http://" + ip + ":8000/" + filename_MD5

	q, _ := qrcode.New(dl_URL, qrcode.Low)
	art := q.ToString(false)
	fmt.Println(art) //输出二维码
	fmt.Println("file:" + file_Path)
	//fmt.Println(filename)
	fmt.Println(dl_URL)
	http.ListenAndServe(":8000", nil)
}
