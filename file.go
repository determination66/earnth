package earnth

import (
	lru "github.com/hashicorp/golang-lru"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

type StaticResourceHandlerOption func(handler *StaticResourceHandler)

type StaticResourceHandler struct {
	dir                     string
	cache                   *lru.Cache
	extensionContentTypeMap map[string]string
	// 大文件不缓存
	maxSize int
}

func NewStaticResourceHandler(dir string, opts ...StaticResourceHandlerOption) (*StaticResourceHandler, error) {
	// 总共缓存 key-value
	c, err := lru.New(1000)
	if err != nil {
		return nil, err
	}
	res := &StaticResourceHandler{
		dir:   dir,
		cache: c,
		// 10 MB，文件大小超过这个值，就不会缓存
		maxSize: 1024 * 1024 * 10,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func StaticWithMaxFileSize(maxSize int) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.maxSize = maxSize
	}
}

func StaticWithCache(c *lru.Cache) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.cache = c
	}
}

func StaticWithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		for ext, contentType := range extMap {
			h.extensionContentTypeMap[ext] = contentType
		}
	}
}

func (s *StaticResourceHandler) Handle(ctx *Context) {
	// 无缓存
	// 1. 拿到目标文件名
	// 2. 定位到目标文件，并且读出来
	// 3. 返回给前端

	// 有缓存
	file, err := ctx.PathValue("file")
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("your request url is invalid")
		return
	}

	dst := filepath.Join(s.dir, file)
	ext := filepath.Ext(dst)[1:]
	header := ctx.Resp.Header()

	if data, ok := s.cache.Get(file); ok {
		contentType := s.extensionContentTypeMap[ext]
		header.Set("Content-Type", contentType)
		header.Set("Content-Length", strconv.Itoa(len(data.([]byte))))
		ctx.RespData = data.([]byte)
		ctx.RespStatusCode = http.StatusOK
		return
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("internal server error")
		return
	}
	// 大文件不缓存
	if len(data) <= s.maxSize {
		s.cache.Add(file, data)
	}
	// 可能的有文本文件，图片，多媒体（视频，音频）
	contentType := s.extensionContentTypeMap[ext]
	header.Set("Content-Type", contentType)
	header.Set("Content-Length", strconv.Itoa(len(data)))
	ctx.RespData = data
	ctx.RespStatusCode = http.StatusOK
}
