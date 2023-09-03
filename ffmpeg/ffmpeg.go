package ffmpeg

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/u2takey/ffmpeg-go"
)

type FFmpeg_Stream struct {
	ffmpeg *ffmpeg_go.Stream
}

var logLevel int = LOGLEVEL_INFO

var Hide_Banner bool = false

type FFmpeg_Args map[string]interface{}

func (args *FFmpeg_Args) AddingOutputMpegMetaData(key string, value string) {
	(*args)[key] = value
}

func InitHookResource() {
	ffmpeg_go.LogCompiledCommand = false
	ffmpeg_go.GlobalCommandOptions = append(ffmpeg_go.GlobalCommandOptions, func(cmd *exec.Cmd) {
		logLevel_str := fmt.Sprintf("%d", logLevel)
		cmd.Args = append(cmd.Args, "-loglevel", logLevel_str)
		fmt.Println("Add additional command parameters:", "-loglevel", logLevel_str)
	})
	ffmpeg_go.GlobalCommandOptions = append(ffmpeg_go.GlobalCommandOptions, func(cmd *exec.Cmd) {
		if Hide_Banner {
			cmd.Args = append(cmd.Args, "-hide_banner")
			fmt.Println("Add additional command parameters:", "-hide_banner")
		}
	})
	ffmpeg_go.GlobalCommandOptions = append(ffmpeg_go.GlobalCommandOptions, func(cmd *exec.Cmd) {
		fmt.Println("Compile complete ffmpeg command parameters:", strings.Join(cmd.Args, " "))
	})
}

func SetFFmpegLogLevel(LogLevel int) {
	switch LogLevel {
	case LOGLEVEL_DEBUG, LOGLEVEL_ERROR, LOGLEVEL_FATAL,
		LOGLEVEL_INFO, LOGLEVEL_PANIC, LOGLEVEL_QUIET,
		LOGLEVEL_TRACE, LOGLEVEL_VERBOSE, LOGLEVEL_WARNING:
		logLevel = LogLevel
	}
}

func SetFFmpegHideBanner(hide bool) {
	Hide_Banner = hide
}

func gen_ffmpeg_go_kwarg(args ...FFmpeg_Args) ffmpeg_go.KwArgs {
	ffmpeg_go_kwargs := ffmpeg_go.KwArgs{}
	for _, arg_value := range args {
		for valueIndex := range arg_value {
			value_raw := arg_value[valueIndex]
			var value string
			switch value_raw.(type) {
			case string:
				value = value_raw.(string)
			case int:
				value = fmt.Sprintf("%d", value_raw.(int))
			case int32:
				value = fmt.Sprintf("%d", value_raw.(int32))
			case int64:
				value = fmt.Sprintf("%d", value_raw.(int64))
			case float32:
				value = fmt.Sprintf("%f", value_raw.(float32))
			case float64:
				value = fmt.Sprintf("%f", value_raw.(float64))
			default:
				panic("运行时异常，接收到未知类型参数!")
			}
			ffmpeg_go_kwargs[valueIndex] = value

		}
	}
	return ffmpeg_go_kwargs
}

func Input(filename string, args ...FFmpeg_Args) *FFmpeg_Stream {
	ffmpeg_go_kwargs := gen_ffmpeg_go_kwarg(args...)
	ffmpeg := &FFmpeg_Stream{}
	ffmpeg.ffmpeg = ffmpeg_go.Input(filename, ffmpeg_go_kwargs)
	return ffmpeg
}

func Concat(inputs []*FFmpeg_Stream, args ...FFmpeg_Args) *FFmpeg_Stream {
	ffmpeg_go_kwargs := gen_ffmpeg_go_kwarg(args...)
	ffmpeg_objs := make([]*ffmpeg_go.Stream, 0)
	for _, f := range inputs {
		ffmpeg_objs = append(ffmpeg_objs, f.ffmpeg)
	}
	ffmpeg := &FFmpeg_Stream{}
	ffmpeg.ffmpeg = ffmpeg_go.Concat(ffmpeg_objs, ffmpeg_go_kwargs)
	return ffmpeg
}

func (ffmpeg *FFmpeg_Stream) Output(filename string, args ...FFmpeg_Args) *FFmpeg_Stream {
	ffmpeg_go_kwargs := gen_ffmpeg_go_kwarg(args...)
	ffmpeg.ffmpeg = ffmpeg.ffmpeg.Output(filename, ffmpeg_go_kwargs)
	return ffmpeg
}

func (ffmpeg *FFmpeg_Stream) Filter(filter_name string, filter_arg ...string) *FFmpeg_Stream {
	var ffmpeg_filter_args ffmpeg_go.Args
	filter_args_raw := make([]string, 0)
	for _, value := range filter_arg {
		filter_args_raw = append(filter_args_raw, value)
	}

	var filter_arg_str string

	for _, arg := range filter_args_raw {
		if len(filter_arg_str) != 0 {
			filter_arg_str = fmt.Sprintf("%s:%s", filter_arg_str, arg)
		} else {
			filter_arg_str = fmt.Sprintf("%s", arg)
		}

	}

	ffmpeg_filter_args = ffmpeg_go.Args{filter_arg_str}

	ffmpeg.ffmpeg = ffmpeg.ffmpeg.Filter(filter_name, ffmpeg_filter_args)
	return ffmpeg
}

func (ffmpeg *FFmpeg_Stream) Run(overToFile bool, verbose bool) error {
	if overToFile {
		ffmpeg.ffmpeg = ffmpeg.ffmpeg.OverWriteOutput()
	}

	if verbose {
		ffmpeg.ffmpeg = ffmpeg.ffmpeg.ErrorToStdOut()
	}

	err := ffmpeg.ffmpeg.Run()
	return err
}

func CheckFFmpegInstalled() (string, bool) {
	cmd := exec.Command("ffmpeg", "-version")

	output, err := cmd.Output()
	if err != nil {
		return "", false
	} else {
		outputline := strings.Split(string(output), "\n")
		version_info := outputline[0]
		re := regexp.MustCompile(`^ffmpeg version (\d.\d).*$`)
		group := re.FindStringSubmatch(version_info)
		if len(group) == 0 {
			return "", false
		} else {
			return group[1], true
		}
	}
}
