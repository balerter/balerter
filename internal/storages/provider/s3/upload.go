package s3

import lua "github.com/yuin/gopher-lua"

func (p *Provider) uploadPNG(L *lua.LState) int {
	p.logger.Debug("CALL S3 UPLOAD PNG")
	return 0
}
