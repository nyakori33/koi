package cqcode

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// 获取文件URL
func GetFileURL(file string) string {
	return fmt.Sprint("file:///", file)
}

func Base64Image(code string) string {
	return "base64://" + code
}

// 查找CQ码
func Find(text string) []string {
	regex := regexp.MustCompile(`(?i)\[CQ:(.|\n)+\]`)
	return regex.FindAllString(text, -1)
}

// 替换CQ码
func Replace(src, repl string) string {
	regex := regexp.MustCompile(`(?i)\[CQ:(.|\n)+\]`)
	return regex.ReplaceAllString(src, repl)
}

// 检查
func Check(code string) bool {
	regex := regexp.MustCompile(`(?i)\[CQ:(.|\n)+\]`)
	return regex.MatchString(code)
}

// 编码CQcode
func Encode(function string, data map[string]interface{}) string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	buffer.WriteString("CQ:")
	buffer.WriteString(function)

	for k, v := range data {
		buffer.WriteString(",")
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(fmt.Sprint(v))
	}

	buffer.WriteString("]")

	return buffer.String()
}

// 解码CQcode
func Decode(code string) (function string, data map[string]string) {
	data = make(map[string]string)

	regex := regexp.MustCompile(`(?i)\[CQ:(.|\n)+\]`)
	code = regex.FindString(code)

	code = strings.TrimPrefix(code, "[")
	code = strings.TrimSuffix(code, "]")

	slices := strings.Split(code, ",")
	regex = regexp.MustCompile(`(?i)CQ:`)
	for _, slice := range slices {
		if regex.MatchString(slice) {
			index := strings.Index(slice, ":")

			function = slice[index+1:]
		}

		if strings.Contains(slice, "=") {
			index := strings.Index(slice, "=")

			data[slice[:index]] = slice[index+1:]
		}
	}

	if strings.EqualFold(function, "json") {
		json := data["data"]

		json = strings.ReplaceAll(json, "&#44;", ",")
		json = strings.ReplaceAll(json, "&amp;", "&")
		json = strings.ReplaceAll(json, "&#91;", "[")
		json = strings.ReplaceAll(json, "&#93;", "]")

		data["data"] = json
	}

	return
}

// 表情
func Face(id int) string {
	data := map[string]interface{}{
		"id": id,
	}

	return Encode("face", data)
}

// 语音
//
// * 可选参数: file, url, magic
func Record(file, url string, magic bool) string {
	data := make(map[string]interface{})

	if file != "" {
		data["file"] = file
	}

	if url != "" {
		data["url"] = url
	}

	if magic {
		data["magic"] = 1
	}

	return Encode("record", data)
}

// 短视频
//
// * 可选参数: cover
func Video(file, cover string) string {
	data := map[string]interface{}{
		"file": file,
	}

	if cover != "" {
		data["cover"] = cover
	}

	return Encode("video", data)
}

// @某人
//
// * 可选参数: name
func At(uid int64, name string) string {
	data := map[string]interface{}{
		"qq": uid,
	}

	if name != "" {
		data["name"] = name
	}

	return Encode("at", data)
}

// 猜拳魔法表情
func Mora() string {
	return "[CQ:rps]"
}

//TODO: 该 CQcode 暂未被 go-cqhttp 支持
// 掷骰子魔法表情
// 窗口抖动(戳一戳)
// 匿名发消息

// 链接分享
//
// * 可选参数: content, image
func LinkShare(url, title, content, image string) string {
	data := map[string]interface{}{
		"url":   url,
		"title": title,
	}

	if content != "" {
		data["content"] = content
	}

	if image != "" {
		data["image"] = image
	}

	return Encode("share", data)
}

//TODO: 该 CQcode 暂未被 go-cqhttp 支持
// 推荐好友/群
// 位置

// 音乐分享
//
// * platform = qq/163/xm
func MusicShare(platform string, id int64) string {
	data := map[string]interface{}{
		"type": platform,
		"id":   id,
	}

	return Encode("music", data)
}

// 音乐自定义分享
//
// * 可选参数: content, image
func MusicShareCustom(url, audio, title, content, image string) string {
	data := map[string]interface{}{
		"type":  "custom",
		"url":   url,
		"audio": audio,
		"title": title,
	}

	if content != "" {
		data["content"] = content
	}

	if image != "" {
		data["image"] = image
	}

	return Encode("music", data)
}

// 图片
//
// * 可选参数: file, url
func Image(file, url string, flash bool) string {
	data := make(map[string]interface{})

	if file != "" {
		data["file"] = file
	}

	if url != "" {
		data["url"] = url
	}

	if flash {
		data["type"] = "flash"
	}

	return Encode("image", data)
}

// 秀图
//
// * 可选参数: file, url
func ShowImage(file, url string, id int) string {
	data := make(map[string]interface{})

	if file != "" {
		data["file"] = file
	}

	if url != "" {
		data["url"] = url
	}

	switch id {
	case 40000, 40001, 40002, 40003, 40004, 40005:
		data["id"] = id
	default:
		data["id"] = 40000
	}

	return Encode("image", data)
}

// 回复
func Reply(message_id int) string {
	data := map[string]interface{}{
		"id": message_id,
	}

	return Encode("reply", data)
}

// 自定义回复
func ReplyCustom(text string, uid, seq, time int64) string {
	data := map[string]interface{}{
		"text": text,
		"qq":   uid,
		"seq":  seq,
		"time": time,
	}

	return Encode("reply", data)
}

// 戳一戳
func Poke(uid int64) string {
	data := map[string]interface{}{
		"qq": uid,
	}

	return Encode("poke", data)
}

// XML消息
func XML(code string) string {
	data := map[string]interface{}{
		"data": code,
	}

	return Encode("xml", data)
}

// JSON消息
func JSON(json string) string {
	json = strings.ReplaceAll(json, ",", "&#44;")
	json = strings.ReplaceAll(json, "&", "&amp;")
	json = strings.ReplaceAll(json, "[", "&#91;")
	json = strings.ReplaceAll(json, "]", "&#93;`")

	data := map[string]interface{}{
		"data": json,
	}

	return Encode("json", data)
}

// CardImage
//
// * 可选参数: source, icon, minwidth, minheight, maxwidth, maxheight
func CardImage(file, source, icon string, minwidth, minheight, maxwidth, maxheight int) string {
	data := map[string]interface{}{
		"file": file,
	}

	if source != "" {
		data["source"] = source
	}

	if icon != "" {
		data["icon"] = icon
	}

	if minwidth != 0 {
		data["minwidth"] = minwidth
	}

	if minheight != 0 {
		data["minheight"] = minheight
	}

	if maxwidth != 0 {
		data["maxwidth"] = maxwidth
	}

	if maxheight != 0 {
		data["maxheight"] = maxheight
	}

	return Encode("cardimage", data)
}

// 文本转语音
func TTS(text string) string {
	data := map[string]interface{}{
		"text": text,
	}

	return Encode("tts", data)
}
