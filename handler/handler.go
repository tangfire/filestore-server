package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// UploadHandler:处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回html页面
		data, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internet server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data,err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		//newFile, err := os.Create("/tmp/" + head.Filename)
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file,err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file,err:%s\n", err.Error())
			return
		}

		//newFile.Seek(0, 0) 的作用是将文件指针移动到文件的开头。这在许多场景中是非常有用的，尤其是在想要重新读取文件的内容时。
		//
		//典型应用场景
		//重新读取文件: 如果你之前已经读取了文件的某些内容，但需要再次从文件的开头开始读取，可以使用 Seek(0, 0) 将指针移动到开头。
		//
		//在写入文件之前清空内容: 在创建文件并写入数据之前，将指针移动到开头是常见的操作。
		//
		//在流式处理中: 在处理大文件或流式数据时，可能会需要多次读取相同的文件，使用 Seek(0, 0) 可以方便地重置指针。
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	}
}

// UploadSucHandler:上传已完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")

}

// GetFileMetaHandler:获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// sha1sum 是一个常见的命令行工具，用于计算和验证文件的 SHA-1 哈希值。
	//SHA-1（Secure Hash Algorithm 1）是一种加密哈希函数，它将任意大小的数据输入转换为固定大小的 160 位（20 字节）哈希值。
	//sha1sum 工具广泛用于文件完整性验证、数据完整性校验以及数字签名。
	// example: sha1sum filename.txt

	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler:下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 因为要在浏览器展现效果，所以需要多这几步操作
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\""+fm.FileName+"\"")

	w.Write(data)

}

// FileMetaUpdateHandler:更新元信息接口(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// FileDeteHandler:删除文件及元信息
func FileDeteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}
