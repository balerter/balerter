package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	lua "github.com/yuin/gopher-lua"
	"hash/crc32"
	"strconv"
	"strings"
	"time"
)

const (
	defaultUploadTimeout = time.Second * 10
)

var (
	crcTable = crc32.MakeTable(crc32.Castagnoli)
)

func (p *Provider) getArgs(L *lua.LState) ([]byte, string, error) {
	dataValue := L.Get(1)
	if dataValue.Type() != lua.LTString {
		return nil, "", fmt.Errorf("upload data must be a string")
	}

	data := []byte(dataValue.String())

	var filename string

	filenameValue := L.Get(2)
	switch filenameValue.Type() {
	case lua.LTNil:
		filename = strconv.Itoa(int(time.Now().UnixNano())) + "-" + strconv.Itoa(int(crc32.Checksum(data, crcTable)))
	case lua.LTString:
		filename = filenameValue.String()
		for _, s := range []string{".png", ".jpg", ".jpeg"} {
			if strings.HasSuffix(filename, s) {
				filename = filename[:len(filename)-len(s)]
			}
		}
	default:
	}

	return data, filename, nil
}

func (p *Provider) uploadPNG(L *lua.LState) int {
	return p.upload(L, "png")
}

func (p *Provider) upload(L *lua.LState, extension string) int {

	data, filename, err := p.getArgs(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("wrong arguments: " + err.Error()))
		return 2
	}

	filename = filename + "." + extension

	creds := credentials.NewStaticCredentials(p.key, p.secret, "")

	cfg := &aws.Config{
		Endpoint:    aws.String(p.endpoint),
		Region:      aws.String(p.region),
		Credentials: creds,
	}

	sess := session.Must(session.NewSession())

	svc := s3.New(sess, cfg)

	ctx, ctxCancel := context.WithTimeout(context.Background(), defaultUploadTimeout)
	defer ctxCancel()

	obj := &s3.PutObjectInput{}

	obj.SetACL("public-read")
	obj.SetBucket(p.bucket)
	obj.SetKey(filename)
	obj.SetBody(bytes.NewReader(data))
	obj.SetContentLength(int64(len(data)))
	obj.SetContentType("image/png")

	_, err = svc.PutObjectWithContext(ctx, obj)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("error upload object: " + err.Error()))
		return 2
	}

	resultFilename := fmt.Sprintf("https://%s.%s/%s", p.bucket, p.endpoint, filename)

	L.Push(lua.LString(resultFilename))

	return 1
}
