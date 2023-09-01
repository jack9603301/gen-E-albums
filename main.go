package main

import (
	"container/list"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/hellflame/argparse"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"os"
	"path/filepath"
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
		return false
	}
	return true
}

func ImageToH265Mpeg(file string, output string, scale string) error {
	file_obj, _ := os.Open(file)
	defer file_obj.Close()
	img, _, err := image.Decode(file_obj)
	if err != nil {
		return err
	}
	b := img.Bounds()
	WidthFromImage := b.Max.X
	HeightFromImage := b.Max.Y
	Direction := WidthFromImage >= HeightFromImage
	Direction = Direction
	return nil
}

func ImageScaleCtr(build_path string, filelists *list.List, scale string, non_interactive bool) {

	if !non_interactive {
		fmt.Println(">>>>>>>>>>>>>>>>请选择分辨率<<<<<<<<<<<<<<<<")
		fmt.Println("1. 3840x2160")
		fmt.Println("2. 1920x1080")
		fmt.Println("3. 2160x3840(竖屏)")
		fmt.Println("4. 31080x1920(竖屏)")
		fmt.Println("5. 自定义分辨率")
	}

	fmt.Println(">>>执行图片预处理程序<<<")
	for file := filelists.Front(); file != nil; file = file.Next() {
		file_filename := file.Value.(string)
		file_ext := filepath.Ext(file_filename)
		file_base := filepath.Base(file_filename[:len(file_filename)-len(file_ext)])
		file_full_without_ext := file_filename[:len(file_filename)-len(file_ext)]
		build_file_full_without_ext := fmt.Sprintf("%s%s", build_path, file_base)
		fmt.Printf("===>>> 正在处理：%s%s\n", file_full_without_ext, file_ext)
		ImageToH265Mpeg(file_full_without_ext+file_ext, build_file_full_without_ext+file_ext, scale)
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
	err := parser.Parse(nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	if !*non_interactive {
		fmt.Println(">>>请注意：进入交互式模式!<<<")
	}

	build_path := fmt.Sprintf("%s%s", *root, "/./build")

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

	ImageScaleCtr(build_path, filelists, *scale, *non_interactive)

}
