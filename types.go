package main

type CodecParam struct {
	video *string
	tag   *string
}

type ArgParam struct {
	rate      int
	codec     CodecParam
	framerate string
	scale     string
	pix_fmt   string
}
