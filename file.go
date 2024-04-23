package earnth

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// MaxUpLoadSize MemorySize byte,you can modify it to limit the size you upload
var MaxUpLoadSize int64 = 32 << 20 // 32MB

type FileUpload struct {
	FileField string

	// consider the conflict file name
	// developer need to set
	DstPathFunc func(*multipart.FileHeader) string
}

func NewFileUpload(fileField string) *FileUpload {
	return &FileUpload{
		FileField: fileField,
		DstPathFunc: func(fileHeader *multipart.FileHeader) string {
			// Default destination path function
			err := os.MkdirAll(filepath.Join("testdata", "uploads"), os.ModePerm)
			if err != nil {
				panic(err)
			}
			return filepath.Join("testdata", "uploads", fileHeader.Filename)
		},
	}
}

func (f *FileUpload) Handle() HandleFunc {
	if f.FileField == "" {
		f.FileField = "file"
	}
	return func(ctx *Context) {
		// limit Space you upload
		err := ctx.Req.ParseMultipartForm(MaxUpLoadSize)
		if err != nil {
			panic(err)
		}

		// Step 1: Read the file content
		file, fileHeader, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("fail to upload file,err:" + err.Error())
			return
		}
		if fileHeader.Size > MaxUpLoadSize {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("fail to upload file,your file size is too large")
			return
		}
		defer file.Close()

		// Step 2: Calculate the target path
		dst := f.DstPathFunc(fileHeader)

		// Step 3: Save the file
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("fail to upload" + err.Error())
			return
		}
		defer dstFile.Close()

		// Step 4: Return response
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("fail to upload,err:" + err.Error())
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("upload success")
	}
}

type FileDownload struct {
	Dir string
}

func NewFileDownload(dir string) *FileDownload {
	return &FileDownload{
		Dir: dir,
	}
}

func (f *FileDownload) Handle() HandleFunc {
	return func(ctx *Context) {
		// 用的是 xxx?file=xxx
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		req = filepath.Clean(req)
		dst := filepath.Join(f.Dir, req)
		// 做一个校验，防止相对路径引起攻击者下载了你的系统文件
		// dst, err = filepath.Abs(dst)
		// if strings.Contains(dst, d.Dir) {
		//
		// }
		fn := filepath.Base(dst)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(ctx.Resp, ctx.Req, dst)
	}
}
