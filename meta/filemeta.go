package meta

import (
	mydb "filestore-server/db"
	"sort"
)

// FileMeta:文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// ByUploadTime 实现 sort.Interface
type ByUploadTime []FileMeta

func (a ByUploadTime) Len() int {
	return len(a)
}

func (a ByUploadTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less 方法用于比较两个 FileMeta 的 UploadAt
func (a ByUploadTime) Less(i, j int) bool {
	// 假设 UploadAt 是一个可直接比较的时间字符串
	return a[i].UploadAt < a[j].UploadAt
}

// UpdateFileMeta:新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB:新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)

}

// GetFileMeta:通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

// RemoveFileMeta:删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
