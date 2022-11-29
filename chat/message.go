package chat

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Dingtalk constants
const (
	// 钉钉机器人发消息地址
	APIEndpoint = "https://oapi.dingtalk.com/robot/send"
)

// Message types
const (
	TypeText       = "text"
	TypeLink       = "link"
	TypeMarkdown   = "markdown"
	TypeActionCard = "actionCard"
	TypeFeedCard   = "feedCard"
)

type BtnOrientationType string

// BtnOrientation
const (
	BtnOrientationV BtnOrientationType = "0" // 按钮竖直排列
	BtnOrientationH BtnOrientationType = "1" // 按钮横向排列
)

// IMessage is any config type that can be sent.
type IMessage interface {
	SetMessageType()
	GetMessageType() string
	MergeMessage(boundary string, m1, m2 IMessage) []IMessage
}

// BaseMessage is base type for all message config types.
type BaseMessage struct {
	MessageType string `json:"msgtype"` // 消息类型
}

// text message config.
type TextConfig struct {
	BaseMessage
	Text TextFieldText `json:"text"`
	At   TextFieldAt   `json:"at"`
}

// text message 'text' field struct.
type TextFieldText struct {
	Content string `json:"content"` // 消息内容。
}

// text message 'at' field struct.
type TextFieldAt struct {
	AtMobiles []string `json:"atMobiles"` // 被@人的手机号。
	AtUserIds []string `json:"atUserIds"` // 被@人的用户userid。
	IsAtAll   bool     `json:"isAtAll"`   // 是否@所有人。
}

// link message config.
type LinkConfig struct {
	BaseMessage
	Link LinkFieldLink `json:"link"`
}

// link message 'link' field struct.
type LinkFieldLink struct {
	Text       string `json:"text"`       // 消息内容。如果太长只会部分展示。
	Title      string `json:"title"`      // 消息标题。
	PicUrl     string `json:"picUrl"`     // 图片URL。
	MessageUrl string `json:"messageUrl"` // 点击消息跳转的URL。
}

// markdown message config.
/*
目前只支持markdown语法的子集，具体支持的元素如下：
	标题
	# 一级标题
	## 二级标题
	### 三级标题
	#### 四级标题
	##### 五级标题
	###### 六级标题

	引用
	> A man who stands for nothing will fall for anything.

	文字加粗、斜体
	**bold**
	*italic*

	链接
	[this is a link](http://name.com)

	图片
	![](http://name.com/pic.jpg)

	无序列表
	- item1
	- item2

	有序列表
	1. item1
	2. item2
*/
type MarkdownConfig struct {
	BaseMessage
	Markdown MarkdownFieldMarkdown `json:"markdown"`
	At       MarkdownFieldAt       `json:"at"`
}

// markdown message 'markdown' field struct.
type MarkdownFieldMarkdown struct {
	Title string `json:"title"` // 首屏会话透出的展示内容。
	Text  string `json:"text"`  // markdown格式的消息。
}

// text message 'at' field struct.
type MarkdownFieldAt struct {
	AtMobiles []string `json:"atMobiles"` // 被@人的手机号。
	AtUserIds []string `json:"atUserIds"` // 被@人的用户userid。
	IsAtAll   bool     `json:"isAtAll"`   // 是否@所有人。
}

// actionCard message config.
type ActionCardConfig struct {
	BaseMessage
	ActionCard ActionCardFieldActionCard `json:"actionCard"`
}

// actionCard message 'actionCard' field struct.
type ActionCardFieldActionCard struct {
	Title          string                `json:"title"`          // 首屏会话透出的展示内容。
	Text           string                `json:"text"`           // markdown格式的消息。
	SingleTitle    string                `json:"singleTitle"`    // 单个按钮的标题。(设置此项和singleURL后，btns无效。)
	SingleUrl      string                `json:"singleURL"`      // 点击singleTitle按钮触发的URL。
	BtnOrientation BtnOrientationType    `json:"btnOrientation"` // 0：按钮竖直排列 1：按钮横向排列
	Btns           []ActionCardFieldBtns `json:"btns"`           // 按钮。
}

// actionCard message 'btns' field struct.
type ActionCardFieldBtns struct {
	Title     string `json:"title"`     // 按钮标题。
	ActionUrl string `json:"actionURL"` // 点击按钮触发的URL。
}

// feedCard message config.
type FeedCardConfig struct {
	BaseMessage
	FeedCard FeedCardFieldFeedCard `json:"feedCard"`
}

// feedCard message 'feedCard' field struct.
type FeedCardFieldFeedCard struct {
	Links []FeedFieldLinks `json:"links"`
}

// feedCard message 'links' field struct.
type FeedFieldLinks struct {
	Title      string `json:"title"`      // 单条信息文本。
	MessageUrl string `json:"messageURL"` // 点击单条信息到跳转链接。
	PicUrl     string `json:"picURL"`     // 单条信息后面图片的URL。
}

// set text message type
func (self *TextConfig) SetMessageType() {
	self.MessageType = TypeText
}

// set link message type
func (self *LinkConfig) SetMessageType() {
	self.MessageType = TypeLink
}

// set markdown message type
func (self *MarkdownConfig) SetMessageType() {
	self.MessageType = TypeMarkdown
}

// set actionCard message type
func (self *ActionCardConfig) SetMessageType() {
	self.MessageType = TypeActionCard
}

// set feedCard message type
func (self *FeedCardConfig) SetMessageType() {
	self.MessageType = TypeFeedCard
}

// get text message type
func (self *BaseMessage) GetMessageType() string {
	return self.MessageType
}

// merge text type message
func (self *TextConfig) MergeMessage(boundary string, im1, im2 IMessage) []IMessage {
	m1, ok := im1.(*TextConfig);
	if !ok {
		return []IMessage{im1, im2}
	}
	m2, ok := im2.(*TextConfig);
	if !ok {
		return []IMessage{im1, im2}
	}

	rm := &TextConfig{
		BaseMessage: BaseMessage{
			MessageType: m1.MessageType,
		},
		Text: TextFieldText{},
		At:   TextFieldAt{},
	}
	if m1.At.IsAtAll || m2.At.IsAtAll {
		rm.At.IsAtAll = true
	}
	rm.Text.Content = fmt.Sprintf("%s\n%s\n%s", m1.Text.Content, boundary, m2.Text.Content)
	rm.At.AtMobiles = append(m1.At.AtMobiles, m2.At.AtMobiles...)
	rm.At.AtUserIds = append(m1.At.AtUserIds, m2.At.AtUserIds...)

	return []IMessage{rm}
}

// merge link type message
func (self *LinkConfig) MergeMessage(boundary string, im1, im2 IMessage) []IMessage {
	// cannot merge link type messages
	return []IMessage{im1, im2}
}

// merge markdown type message
func (self *MarkdownConfig) MergeMessage(boundary string, im1, im2 IMessage) []IMessage {
	m1, ok := im1.(*MarkdownConfig);
	if !ok {
		return []IMessage{im1, im2}
	}
	m2, ok := im2.(*MarkdownConfig);
	if !ok {
		return []IMessage{im1, im2}
	}

	rm := &MarkdownConfig{
		BaseMessage: BaseMessage{
			MessageType: m1.MessageType,
		},
		Markdown: MarkdownFieldMarkdown{},
		At:       MarkdownFieldAt{},
	}
	if m1.At.IsAtAll || m2.At.IsAtAll {
		rm.At.IsAtAll = true
	}
	rm.Markdown.Title = fmt.Sprintf("%s|%s", m1.Markdown.Title, m2.Markdown.Title)
	rm.Markdown.Text = fmt.Sprintf(`%s

%s

%s`, m1.Markdown.Text, boundary, m2.Markdown.Text)
	rm.At.AtMobiles = append(m1.At.AtMobiles, m2.At.AtMobiles...)
	rm.At.AtUserIds = append(m1.At.AtUserIds, m2.At.AtUserIds...)

	return []IMessage{rm}
}

// merge actionCard type message
func (self *ActionCardConfig) MergeMessage(boundary string, im1, im2 IMessage) []IMessage {
	// cannot merge actionCard type messages
	return []IMessage{im1, im2}
}

// merge feedCard type message
func (self *FeedCardConfig) MergeMessage(boundary string, im1, im2 IMessage) []IMessage {
	m1, ok := im1.(*FeedCardConfig);
	if !ok {
		return []IMessage{im1, im2}
	}
	m2, ok := im2.(*FeedCardConfig);
	if !ok {
		return []IMessage{im1, im2}
	}

	rm := &FeedCardConfig{
		BaseMessage: BaseMessage{},
		FeedCard: FeedCardFieldFeedCard{
			Links: append(m1.FeedCard.Links, m2.FeedCard.Links...),
		},
	}

	return []IMessage{rm}
}

// unmarshal IMessage bytes.
func UnmarshalBytes(bytes []byte) (im IMessage, err error) {
	base := BaseMessage{}
	err = json.Unmarshal(bytes, &base)
	if err != nil {
		return
	}

	msgType := base.GetMessageType()
	switch msgType {
	case TypeText:
		config := &TextConfig{}
		err = json.Unmarshal(bytes, config)
		if err != nil {
			return
		}
		im = config
		return
	case TypeLink:
		config := &LinkConfig{}
		err = json.Unmarshal(bytes, config)
		if err != nil {
			return
		}
		im = config
		return

	case TypeMarkdown:
		config := &MarkdownConfig{}
		err = json.Unmarshal(bytes, config)
		if err != nil {
			return
		}
		im = config
		return
	case TypeActionCard:
		config := &ActionCardConfig{}
		err = json.Unmarshal(bytes, config)
		if err != nil {
			return
		}
		im = config
		return
	case TypeFeedCard:
		config := &FeedCardConfig{}
		err = json.Unmarshal(bytes, config)
		if err != nil {
			return
		}
		im = config
		return
	}
	return nil, errors.New("Failed to unmarshal message.")
}
