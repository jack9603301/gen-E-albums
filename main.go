package main

import (
	"container/list"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/h2non/filetype"
	"github.com/hellflame/argparse"
	_ "github.com/u2takey/ffmpeg-go"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var filelists = list.New()

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

func ImageToH265Mpeg(file string, output string, framerate float64) error {
	fmt.Println("转入图片压制视频处理程序，同时应用滤镜")
	fmt.Println("输入图片的采样帧率是：", framerate)
	//err := ffmpeg.Input(file, ffmpeg.KwArgs{"framerate": framerate})
	return nil
}

func ImageToScaleImage(file string, output string, scale string) error {
	file_obj, _ := os.Open(file)
	defer file_obj.Close()
	img, _, err := image.Decode(file_obj)
	if err != nil {
		return err
	}
	b := img.Bounds()
	var WidthFromImage uint64 = uint64(b.Max.X)
	var HeightFromImage uint64 = uint64(b.Max.Y)
	Direction := WidthFromImage >= HeightFromImage
	scale_slice := strings.Split(scale, "x")
	if len(scale_slice) < 2 {
		scale_slice = strings.Split(scale, "X")
	}
	Direction = Direction
	var WidthToImage uint64 = 0
	var HeightToImage uint64 = 0
	var Scale uint64
	if Direction {
		WidthToImage, _ = strconv.ParseUint(scale_slice[0], 10, 0)
		switch {
		case WidthToImage > WidthFromImage:
			Scale = WidthToImage / WidthFromImage
		default:
			Scale = WidthFromImage / WidthToImage
		}
		HeightToImage = HeightFromImage * Scale
	} else {
		HeightToImage, _ = strconv.ParseUint(scale_slice[1], 10, 0)
		switch {
		case HeightToImage > HeightFromImage:
			Scale = HeightToImage / HeightFromImage
		default:
			Scale = HeightToImage / HeightFromImage
		}
		WidthToImage = WidthFromImage * Scale
	}
	fmt.Println("======请注意：使用缩放比率弹性选择")
	fmt.Println("选定线性缩放参数：", scale)

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err = mw.ReadImage(file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("图片预处理原始大小为: %dx%d\n", mw.GetImageWidth(), mw.GetImageHeight())

	err = mw.AdaptiveResizeImage(uint(WidthToImage), uint(HeightToImage))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("图片预处理后的实际大小为: %dx%d\n", mw.GetImageWidth(), mw.GetImageHeight())

	fmt.Println("输出文件为：", output)

	err = mw.WriteImage(output)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func ImageScaleCtr(build_path string, filelists *list.List, scale string, non_interactive bool, framerate float64) {

	if !non_interactive {
		width, height := getResolution()

		if width == 0 {
			// 输出一个空行
			fmt.Println()
			width, height = requestInputResolution()
		}

		scale = fmt.Sprintf("%dx%d", width, height)

	}

	fmt.Println(">>>执行图片预处理程序<<<")
	for file := filelists.Front(); file != nil; file = file.Next() {
		file_filename := file.Value.(string)
		file_ext := filepath.Ext(file_filename)
		file_base := filepath.Base(file_filename[:len(file_filename)-len(file_ext)])
		file_full_without_ext := file_filename[:len(file_filename)-len(file_ext)]
		build_file_full_without_ext := fmt.Sprintf("%s%s", build_path, file_base)
		fmt.Printf("===>>> 正在处理：%s%s\n", file_full_without_ext, file_ext)
		ImageToScaleImage(file_full_without_ext+file_ext, build_file_full_without_ext+file_ext, scale)
		mpeg_file_ext := ".mp4"
		ImageToH265Mpeg(build_file_full_without_ext+file_ext, build_file_full_without_ext+mpeg_file_ext, framerate)
	}
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
	framerate_str := parser.String("fr", "framerate", &argparse.Option{
		Help:     "Input picture sampling frame rate",
		Required: false,
		Default:  "1/2"})
	err := parser.Parse(nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	framerate_pragam, err := expr.Compile(*framerate_str, expr.Env(nil))
	if err != nil {
		fmt.Println(err)
		return
	}

	framerate, err := expr.Run(framerate_pragam, expr.Env(nil))
	if err != nil {
		fmt.Println(err)
		return
	}

	if !*non_interactive {
		fmt.Println(">>>请注意：进入交互式模式!<<<")
	}

	imagick.Initialize()
	defer imagick.Terminate()

	build_path := fmt.Sprintf("%s%s", *root, "/./build/")

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

	ImageScaleCtr(build_path, filelists, *scale, *non_interactive, framerate.(float64))

}
