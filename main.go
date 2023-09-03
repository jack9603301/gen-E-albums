package main

import (
	"container/list"
	"fmt"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gen-E-albums/ffmpeg"
	"github.com/h2non/filetype"
	"github.com/hellflame/argparse"
	"gopkg.in/gographics/imagick.v3/imagick"

	"gen-E-albums/assert"
)

var filelists = list.New()
var MpegHeight uint64
var MpegWidth uint64

const DESCMeta string = "description=\"此视频相册由gen-E-albums电子相册编译程序压制而成，它是出于兴趣和需求而编写，源代码从github的jack9603301/gen-E-albums获取\""

func searchFile(path string, info os.FileInfo, err error) error {
	kind, _ := filetype.MatchFile(path)
	fmt.Println(path, "=>", kind.MIME.Value)
	if strings.HasPrefix(kind.MIME.Value, "image/") {
		filelists.PushBack(path)
	}
	return nil
}

func checkexist_buildpath(build string) bool {
	_, err := os.Stat(build)
	if err == nil {
		return true
	}
	return false
}

func ImageToH265Mpeg(file string, output string, args_param ArgParam) (string, error) {
	fmt.Println("转入图片压制视频处理程序，同时应用滤镜")
	fmt.Println("输入图片的保留时间是：", args_param.duration)
	fmt.Println("输入图片路径：", file)
	codec := args_param.codec
	ffmpeg_output_KwArg := ffmpeg.FFmpeg_Args{}
	assert.Assert(codec.video != nil, "参数错误，没有编码器信息！")
	if codec.video == nil {
		panic("参数错误，没有编码器信息！")
	}
	ffmpeg_output_KwArg["c:v"] = *codec.video
	if *codec.video == string("libx265") {
		tag := codec.tag
		if tag != nil {
			ffmpeg_output_KwArg["tag:v"] = *tag
		}
	}
	ffmpeg_output_KwArg["r"] = args_param.rate
	ffmpeg_output_KwArg["pix_fmt"] = args_param.pix_fmt
	ffmpeg_output_KwArg["t"] = args_param.duration
	ffmpeg_output_KwArg.AddingOutputMpegMetaData("metadata", DESCMeta)
	ffmpeg_input_KwArg := ffmpeg.FFmpeg_Args{}
	ffmpeg_input_KwArg["loop"] = 1

	stream := ffmpeg.Input(file, ffmpeg_input_KwArg).
		Filter("fade", "t=in", "st=0", "d=0.5").
		Filter("fade", "t=out", "st=1.5", "d=0.5").
		Output(output, ffmpeg_output_KwArg)

	err := stream.Run(true, true)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("单图片视频临时文件输出至：", output)
	return output, nil
}

func ImageToScaleImage(file string, output string, scale string) error {
	file_obj, _ := os.Open(file)
	defer file_obj.Close()
	img, _, err := image.Decode(file_obj)
	if err != nil {
		return err
	}
	b := img.Bounds()
	scale_slice := strings.Split(scale, "x")
	if len(scale_slice) < 2 {
		scale_slice = strings.Split(scale, "X")
	}
	var WidthFromImage uint64 = uint64(b.Max.X)
	var HeightFromImage uint64 = uint64(b.Max.Y)
	ParamWidth, _ := strconv.ParseUint(scale_slice[0], 10, 0)
	ParamHeight, _ := strconv.ParseUint(scale_slice[1], 10, 0)
	WidthFlag := WidthFromImage >= ParamWidth
	HeightFlag := HeightFromImage >= ParamHeight
	var Direction bool = false
	switch {
	case WidthFlag == true && HeightFlag == true:
		Direction = WidthFromImage < HeightFromImage
	case WidthFlag == true && HeightFlag == false:
		Direction = true
	case WidthFlag == true && HeightFlag == false:
		Direction = false
	}
	var WidthToImage uint64 = 0
	var HeightToImage uint64 = 0
	if Direction {
		WidthToImage = ParamWidth
		HeightToImage = HeightFromImage * WidthToImage / WidthFromImage
	} else {
		HeightToImage = ParamHeight
		WidthToImage = WidthFromImage * HeightToImage / HeightFromImage
	}
	fmt.Println("======请注意：使用缩放比率弹性选择")
	fmt.Printf("选定线性缩放参数： %dx%d\n", WidthToImage, HeightToImage)

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	fmt.Println("创建背景调色板")
	bg := imagick.NewPixelWand()
	bg.SetColor("none")
	bg.SetAlpha(0)

	fmt.Println("创建背景图层")
	mw.NewImage(uint(ParamWidth), uint(ParamHeight), bg)
	mw.SetImageFormat("png")

	fmt.Printf("读取 %s 图片文件\n", file)
	source_mw := imagick.NewMagickWand()
	defer source_mw.Destroy()
	source_mw.ReadImage(file)

	fmt.Printf("图片预处理原始大小为: %dx%d\n", source_mw.GetImageWidth(), source_mw.GetImageHeight())

	err = source_mw.AdaptiveResizeImage(uint(WidthToImage), uint(HeightToImage))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("图片缩放处理后的大小为: %dx%d\n", source_mw.GetImageWidth(), source_mw.GetImageHeight())

	fmt.Println("将实际图片文件覆盖到背景图片中")
	offsetX := int((uint(ParamWidth) - uint(source_mw.GetImageWidth())) / 2)
	offsetY := int((uint(ParamHeight) - uint(source_mw.GetImageHeight())) / 2)
	mw.CompositeImage(source_mw, imagick.COMPOSITE_OP_COPY, true, offsetX, offsetY)

	fmt.Println("图片处理完毕，写入至:", output)
	mw.WriteImage(output)

	return nil
}

func ImageScaleCtr(build_path string, filelists *list.List, non_interactive bool, args_param ArgParam) (*list.List, error) {

	var scale string

	if !non_interactive {
		width, height := getResolution()

		if width == 0 {
			// 输出一个空行
			fmt.Println()
			width, height = requestInputResolution()
		}

		scale = fmt.Sprintf("%dx%d", width, height)

	} else {
		scale = args_param.scale
	}

	fmt.Println(">>>执行图片预处理程序<<<")
	outputMpegImages := list.New()
	for file := filelists.Front(); file != nil; file = file.Next() {
		file_filename := file.Value.(string)
		file_ext := filepath.Ext(file_filename)
		file_base := filepath.Base(file_filename[:len(file_filename)-len(file_ext)])
		file_full_without_ext := file_filename[:len(file_filename)-len(file_ext)]
		build_file_full_without_ext := fmt.Sprintf("%s%s", build_path, file_base)
		fmt.Printf("===>>> 正在处理：%s%s\n", file_full_without_ext, file_ext)
		png_ext := ".png"
		ImageToScaleImage(file_full_without_ext+file_ext, build_file_full_without_ext+png_ext, scale)
		mpeg_file_ext := ".mp4"
		output_filename, err := ImageToH265Mpeg(build_file_full_without_ext+png_ext, build_file_full_without_ext+mpeg_file_ext, args_param)
		if err != nil {
			fmt.Println(err)
			return outputMpegImages, err
		}
		outputMpegImages.PushBack(output_filename)
	}
	return outputMpegImages, nil
}

func VideoConcat(ImagesMpeg *list.List, output string, args_param ArgParam) error {
	var ImageObjs []*ffmpeg.FFmpeg_Stream
	for file := ImagesMpeg.Front(); file != nil; file = file.Next() {
		ImageObj := ffmpeg.Input(file.Value.(string))
		ImageObjs = append(ImageObjs, ImageObj)
	}

	codec := args_param.codec
	ffmpeg_output_KwArg := ffmpeg.FFmpeg_Args{}
	assert.Assert(codec.video != nil, "参数错误，没有编码器信息！")
	if codec.video == nil {
		panic("参数错误，没有编码器信息！")
	}
	ffmpeg_output_KwArg["c:v"] = *codec.video
	if *codec.video == string("libx265") {
		tag := codec.tag
		if tag != nil {
			ffmpeg_output_KwArg["tag:v"] = *tag
		}
	}
	ffmpeg_output_KwArg["r"] = args_param.rate

	stream := ffmpeg.Concat(ImageObjs, ffmpeg.FFmpeg_Args{}).
		Output(output, ffmpeg_output_KwArg)

	err := stream.Run(true, true)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	parser := argparse.NewParser("gen-E-albums", "电子视频相册编译程序帮助", nil)
	root := parser.String("r", "root", &argparse.Option{
		Help:       "Source image file directory",
		Positional: true})
	scale := parser.String("s", "scale", &argparse.Option{
		Help:     "Output image scaling control",
		Required: false,
		Default:  "1920x1080"})
	non_interactive := parser.Flag("ni", "non-interactive", &argparse.Option{
		Help:     "Enable non-interactive mode",
		Required: false})
	duration := parser.Int("d", "duration", &argparse.Option{
		Help:     "The picture shows the duration",
		Required: true,
		Default:  "2"})
	rate := parser.Int("r", "rate", &argparse.Option{
		Help:     "Output video frame rate",
		Required: true,
		Default:  "25"})
	pix_fmt := parser.String("pf", "pix_fmt", &argparse.Option{
		Help:     "original image format",
		Required: true,
		Default:  "yuv444p10le"})
	vcodec_param := parser.String("vc", "vcodec", &argparse.Option{
		Help:     "Video encoder settings",
		Required: true,
		Default:  "libx265"})
	tag_param := parser.String("t", "tag", &argparse.Option{
		Help:     "Video encoder tag (Only for H265)",
		Required: false,
		Default:  "hvc1"})
	log_level := parser.String("ll", "log_level", &argparse.Option{
		Help:     "Set Log level",
		Required: false,
		Default:  "INFO"})
	err := parser.Parse(nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch {
	case len(*root) == 0:
		fmt.Println("错误，参数错误，ROOT不能传入空参数")
		return
	case len(*scale) == 0:
		fmt.Println("错误，参数错误，SCALE不能传入空参数")
		return
	case len(*pix_fmt) == 0:
		fmt.Println("错误，参数错误，PIX_FMT不能传入空参数")
		return
	case len(*vcodec_param) == 0:
		fmt.Println("错误，参数错误，VCODEC_PARAM不能传入空参数")
		return
	case len(*log_level) == 0:
		fmt.Println("错误，参数错误，LOG_LEVEL不能传入空参数")
		return
	}

	log_level_Upper := strings.ToUpper(*log_level)
	log_level_parse := ffmpeg.LOGLEVEL_INFO

	switch log_level_Upper {
	case "QUIET":
		log_level_parse = ffmpeg.LOGLEVEL_QUIET
	case "PANIC":
		log_level_parse = ffmpeg.LOGLEVEL_PANIC
	case "FATAL":
		log_level_parse = ffmpeg.LOGLEVEL_FATAL
	case "ERROR":
		log_level_parse = ffmpeg.LOGLEVEL_ERROR
	case "WARNING":
		log_level_parse = ffmpeg.LOGLEVEL_WARNING
	case "INFO":
		log_level_parse = ffmpeg.LOGLEVEL_INFO
	case "VERBOSE":
		log_level_parse = ffmpeg.LOGLEVEL_VERBOSE
	case "DEBUG":
		log_level_parse = ffmpeg.LOGLEVEL_DEBUG
	case "TRACE":
		log_level_parse = ffmpeg.LOGLEVEL_TRACE
	default:
		fmt.Println(">>> 请注意，日志类型", log_level_Upper, "无效！自动选择INFO级别！")
		log_level_parse = ffmpeg.LOGLEVEL_INFO
	}

	version, installed := ffmpeg.CheckFFmpegInstalled()

	if !installed {
		fmt.Println(">>>无法检测到ffmpeg工具，请安装ffmpeg并加入环境变量！")
		return
	} else {
		fmt.Println(">>> 检测到ffmpeg版本：", version)
	}
	ffmpeg.SetFFmpegLogLevel(log_level_parse)
	ffmpeg.SetFFmpegHideBanner(true)

	fmt.Println(">>>初始化FFmpeg资源！")
	ffmpeg.InitHookResource()

	vcodec := *vcodec_param
	tag := *tag_param

	codec := CodecParam{
		video: &vcodec,
		tag:   &tag,
	}

	if len(tag) == 0 || *vcodec_param != "libx265" {
		codec.tag = nil
	}

	args_param := ArgParam{
		rate:     *rate,
		duration: *duration,
		codec:    codec,
		scale:    *scale,
		pix_fmt:  *pix_fmt,
	}

	if !*non_interactive {
		fmt.Println(">>>请注意：进入交互式模式!<<<")
	}

	imagick.Initialize()
	defer imagick.Terminate()

	root_dir := strings.TrimSuffix(*root, "/")

	build_path := fmt.Sprintf("%s%s", root_dir, "/build/")

	if !checkexist_buildpath(build_path) {
		os.Mkdir(build_path, os.ModePerm)
	}

	fmt.Println(">>>开始检测文件类型<<<")
	err = filepath.Walk(*root, searchFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(">>>检测文件类型完成<<<")
	fmt.Println(">>>以下文件进入处理程序<<<")

	for file := filelists.Front(); file != nil; file = file.Next() {
		file_filename := file.Value.(string)
		file_ext := filepath.Ext(file_filename)
		file_base := file_filename[:len(file_filename)-len(file_ext)]
		fmt.Printf("%s%s\n", file_base, file_ext)
	}

	outputMpegImages, _ := ImageScaleCtr(build_path, filelists, *non_interactive, args_param)
	if err != nil {
		fmt.Println(err)
		return
	}

	output_mp4 := fmt.Sprintf("%s%s", build_path, "output.mp4")

	VideoConcat(outputMpegImages, output_mp4, args_param)

}
