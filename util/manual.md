util.go定义了一个处理文件元信息的模块，包含了文件的基本信息结构体、相关的操作函数，以及对这些信息的排序功能。

### 1. 包声明与导入
```go
package meta

import (
	mydb "filestore-server/db"
	"sort"
)
```
- **`package meta`**: 定义了一个名为 `meta` 的包。
- **`import`**: 导入了其他包，`mydb` 用于数据库操作，`sort` 用于排序。

### 2. 文件元信息结构体
```go
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}
```
- **`FileMeta`**: 结构体，存储文件的元信息，包括：
    - `FileSha1`: 文件的 SHA-1 哈希值。
    - `FileName`: 文件名称。
    - `FileSize`: 文件大小（字节）。
    - `Location`: 文件存储位置。
    - `UploadAt`: 文件上传时间。

### 3. 文件元信息映射
```go
var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}
```
- **`fileMetas`**: 一个映射（map），以 SHA-1 值作为键，`FileMeta` 结构体作为值。
- **`init()`**: 初始化函数，创建 `fileMetas` 映射。

### 4. 排序功能
```go
type ByUploadTime []FileMeta

func (a ByUploadTime) Len() int {
	return len(a)
}

func (a ByUploadTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByUploadTime) Less(i, j int) bool {
	return a[i].UploadAt < a[j].UploadAt
}
```
- **`ByUploadTime`**: 自定义类型，实现 `sort.Interface`，用于按上传时间排序。
- **`Len()`**: 返回切片的长度。
- **`Swap()`**: 交换切片中的两个元素。
- **`Less()`**: 比较两个 `FileMeta` 的 `UploadAt` 字段，以决定它们的顺序。

### 5. 文件元信息更新和获取
```go
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}
```
- **`UpdateFileMeta(fmeta FileMeta)`**: 更新内存中的文件元信息。
- **`UpdateFileMetaDB(fmeta FileMeta)`**: 将文件元信息保存到数据库，返回操作结果。
- **`GetFileMeta(fileSha1 string)`**: 根据 SHA-1 值获取内存中的文件元信息。

### 6. 从数据库获取文件元信息
```go
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}

	return fmeta, nil
}
```
- **`GetFileMetaDB(fileSha1 string)`**: 从数据库中获取文件元信息，返回 `FileMeta` 结构体和可能的错误。

### 7. 获取最新文件元信息
```go
func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}
```
- **`GetLastFileMetas(count int)`**: 获取最近的文件元信息，返回最新的 `count` 个文件的元信息。

### 8. 删除文件元信息
```go
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
```
- **`RemoveFileMeta(fileSha1 string)`**: 根据 SHA-1 值删除对应的文件元信息。

### 总结
util.go实现了文件元信息的存储、更新、获取和排序功能，支持在内存和数据库之间进行操作。通过 `FileMeta` 结构体来描述文件的基本信息，并提供了便捷的方法来管理这些信息。