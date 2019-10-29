package main

type Storage interface {
	// url，过期时间
	// url转为短地址
	Shorten(url string, exp int64) (string, error)
	// 短地址信息
	ShortLinkInfo(short string) (interface{}, error)
	// 短地址转回长地址
	Unshort(lang string) (string, error)
}
