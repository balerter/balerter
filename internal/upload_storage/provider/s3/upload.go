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

func (p *Provider) getArgs(luaState *lua.LState) (data []byte, filename string, err error) {
	dataValue := luaState.Get(1)
	if dataValue.Type() != lua.LTString {
		return nil, "", fmt.Errorf("upload data must be a string")
	}

	data = []byte(dataValue.String())

	filenameValue := luaState.Get(2) //nolint:gomnd // param position
	switch filenameValue.Type() {
	case lua.LTNil:
		filename = strconv.Itoa(int(time.Now().UnixNano())) + "-" + strconv.Itoa(int(crc32.Checksum(data, crcTable)))
	case lua.LTString:
		filename = filenameValue.String()
		for _, s := range []string{".png", ".jpg", ".jpeg"} {
			filename = strings.TrimSuffix(filename, s)
		}
	default:
	}

	return data, filename, nil
}

func (p *Provider) uploadPNG(luaState *lua.LState) int {
	return p.upload(luaState, "png")
}

func (p *Provider) upload(luaState *lua.LState, extension string) int {
	data, filename, err := p.getArgs(luaState)
	if err != nil {
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("wrong arguments: " + err.Error()))
		return 2 //nolint:gomnd // params count
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
		luaState.Push(lua.LNil)
		luaState.Push(lua.LString("error upload object: " + err.Error()))
		return 2 //nolint:gomnd // params count
	}

	resultFilename := fmt.Sprintf("https://%s.%s/%s", p.bucket, p.endpoint, filename)

	luaState.Push(lua.LString(resultFilename))

	return 1
}
