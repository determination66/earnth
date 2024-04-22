package earnth

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

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
		// Step 1: Read the file content
		file, fileHeader, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("fail to upload file")
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
