package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func requestInputResolution() (int, int) {
	for {
		inputStr := requestInput("请输入你需要的分辨率: ", true)
		// 处理用户输入了的情况，这里要用正则判断输入格式
		re := regexp.MustCompile(`\s{0,}(\d+)\s{0,}x\s{0,}(\d+)\s{0,}`)
		matchGroup := re.FindStringSubmatch(inputStr)
		if len(matchGroup) >= 3 {
			width, _ := strconv.Atoi(matchGroup[1])
			height, _ := strconv.Atoi(matchGroup[2])
			return width, height
		}
	}
}

// 请求用户输入，选择节点
func requestInput(msg string, force bool) string {
	var inputStr string
	fmt.Print(msg)
	for {
		// fmt.Scanln() 用于从扫描用户的输入，如果输入的字符串左右有空格，
		// fmt.Scanln()会自动去除，所以无需用strings.TrimSapce()函数去除空格
		if _, err := fmt.Scanln(&inputStr); err == nil {
			// 暂时用个变量假设用户输入格式正确
			if inputStr != "" {
				return inputStr
			}
		}
		if !force {
			break
		}
		fmt.Print("\033[1A\033[K")
		fmt.Print(msg)
	}
	return ""
}

// selectOptions 用于请求用户输入想要选择的选项对应的序号
func selectOptions(options []string, tips map[string]string, defaultNum string) string {
	header := tips["header"]
	if tips["useErr"] == "1" {
		header = tips["err"]
	}

	footer := tips["footer"]
	if footer == "" {
		footer = "请输入{selectRange}其中一个数字并回车确定(默认" + defaultNum + "): "
	}

	selectRange := "[1-" + strconv.Itoa(len(options)) + "]"
	footer = strings.Replace(footer, "{selectRange}", selectRange, 1)

	numStrToSelect := ""
	body := ""
	for i, option := range options {
		num := i + 1
		numStrToSelect += strconv.Itoa(num) + ","
		option = "- " + strconv.Itoa(num) + "." + option + "\n"
		body += option
	}
	numStrToSelect = strings.TrimLeft(numStrToSelect, ",")

	//---------------- 输出到终端的文字 开始 --------------------
	header = strings.TrimLeft(header, "\n")
	tip := "\n" + header + "\n" + body + footer

	//---------------- 向输入流请求输入参数 开始 --------------------
	//向输入流请求输入参数
	input := requestInput(tip, false)
	if input == "" {
		input = defaultNum
	}
	fmt.Println("您选择的是: ", input)

	if !strings.Contains(numStrToSelect, input) {
		//输出这个可以用来清屏
		// echo "\033c";
		tips["useErr"] = "1"
		err := ""
		defaultErr := "您的输入为“{input}”，不在[y/n]范围，请重新输入(y/n, 默认" + defaultNum + "): "
		if tips["errTemplate"] == "" {
			err = defaultErr
		} else {
			err = tips["errTemplate"]
		}
		err = strings.Replace(err, "{input}", input, 1)
		// err = strings.TrimLeft(err, "\n")
		tips["err"] = strings.Replace(err, "{selectRange}", selectRange, 1)
		input = selectOptions(options, tips, defaultNum)
		return input
	}

	return input
}

// getResolution 用于根据用户输入的序号，获取序号对应的分辨率
func getResolution() (int, int) {
	defaultNum := "2"
	options := []string{
		"3840x2160",
		"1920x1080",
		"2160x3840(竖屏)",
		"1080x1920(竖屏)",
		"自定义宽高",
	}
	tips := map[string]string{
		"header":      "请选择要生成的视频分辨率，请根据图片大小选择，如果图片比较不高清，不建议选3840x2160: ",
		"footer":      `请输入{selectRange}其中一个数字并回车确定(默认` + defaultNum + `): `,
		"errTemplate": "您的输入为{input}，不在{selectRange}范围，请重新输入(默认" + defaultNum + "): ",
		"useErr":      "0",
	}
	numStr := selectOptions(options, tips, defaultNum)
	num, _ := strconv.Atoi(numStr)
	num -= 1

	resolutionArr := [][]int{
		{3840, 2160},
		{1920, 1080},
		{2160, 3840},
		{1080, 1920},
		{0, 0},
	}
	// frameRateNum := getFrameRateNum("")
	resolution := resolutionArr[num]
	return resolution[0], resolution[1]
}
