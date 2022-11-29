package chat

import (
	"crypto/tls"
	"dingtalk-push/conf"
	loggerSystem "dingtalk-push/utils/log/system"
	"net/http"
	"sync"
	"time"
)

// 机器人池， 控制消息频率
type BotPool struct {
	Bots            map[string]*BotAPI // 机器人   secret: BotAPI
	lock            sync.Mutex
	MessageDuration int    // 消息间隔
	Boundary        string // 多条消息分界线
	chMessage       chan MessageChan
}

// 消息通道
type MessageChan struct {
	Secret      string
	AccessToken string
	Message     IMessage
}

var BotPoolInstance *BotPool
var onceBotPool sync.Once
var client *http.Client
var onceClient sync.Once

func init() {
	onceBotPool.Do(initBotPool)
	onceClient.Do(initClient)
	go BotPoolInstance.limiter()
}

// 初始化BotPool
func initBotPool() {
	BotPoolInstance = &BotPool{
		Bots:            make(map[string]*BotAPI),
		lock:            sync.Mutex{},
		MessageDuration: conf.ConfigYamlInstance.DingtalkConfig.MessageDuration,
		Boundary:        conf.ConfigYamlInstance.DingtalkConfig.Boundary,
		chMessage:       make(chan MessageChan, 1000),
	}
}

// 初始化http client
func initClient() {
	timeout := time.Duration(conf.ConfigYamlInstance.DingtalkConfig.TimeOut) * time.Second
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
		Timeout: timeout,
	}
}

// 判断机器人是否存在
func (self *BotPool) isExist(secret string) bool {
	self.lock.Lock()
	defer self.lock.Unlock()

	_, ok := self.Bots["secret"]
	return ok
}

// 添加机器人到机器人池
func (self *BotPool) addBot(secret, accessToken string) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	botApi, err := NewBotAPIWithClient(secret, accessToken, client)
	if err != nil {
		return err
	}
	self.Bots[secret] = botApi
	return nil
}

// 发送消息
func (self *BotPool) Send(secret, accessToken string, m IMessage) bool {
	select {
	case self.chMessage <- MessageChan{
		Secret:      secret,
		AccessToken: accessToken,
		Message:     m,
	}:
		return true
	default:
	}
	return false
}

// 消息限频
// 每个频率区间遍历一次消息通道。 同一个机器人在每个频率区间只会发送一条消息（合并后的消息算一条），多余的消息会回写回消息通道
func (self *BotPool) limiter() {
	for {
		// 无法合并的消息，重新写进chMessage，下一次频率发送
		otherMessage := make([]MessageChan, 0)
		if len(self.chMessage) > 0 {
			messageMap := make(map[string]map[string]IMessage) // secret: MessageType: IMessage
			for message := range self.chMessage {
				secret := message.Secret
				accessToken := message.AccessToken
				m := message.Message
				msgType := m.GetMessageType()

				// 判断机器人是否在池中
				if !self.isExist(secret) {
					if err := self.addBot(secret, accessToken); err != nil {
						continue
					}
				}

				// 判断机器人
				mtMap, ok := messageMap[secret]
				if !ok {
					// 机器人不存在，则添加机器人
					mtMap = map[string]IMessage{
						msgType: m,
					}
					messageMap[secret] = mtMap
				} else {
					// 机器人存在，判断消息类型
					_, ok := mtMap[msgType]
					if !ok {
						// 若消息类型不存在，则回写消息，下个频率周期发送
						otherMessage = append(otherMessage, message)
					} else {
						// 机器人和消息类型都存在， 则合并消息
						// rim1:初始消息合并后的消息，  rim2:无法合并的消息
						rim1, rim2 := self.mergeMessage(mtMap[msgType], m)
						messageMap[secret][msgType] = rim1
						// 回写无法合并的消息
						for _, rim := range rim2 {
							otherMessage = append(otherMessage, MessageChan{
								Secret:      secret,
								AccessToken: accessToken,
								Message:     rim,
							})
						}
					}
				}

				if len(self.chMessage) == 0 {
					break
				}
			}

			// 回写无法合并的消息
			for _, im := range otherMessage {
				self.chMessage <- im
			}

			// 发消息
			for secret, mtMap := range messageMap {
				bot := self.Bots[secret]
				go func(secret string, mtMap map[string]IMessage) {
					for _, m := range mtMap {
						ok, err := bot.Send(m)
						if !ok || err != nil {
							loggerSystem.ErrorWithFields(
								"send error",
								loggerSystem.Fields{
									"message": m,
									"secret":  secret,
									"error":   err,
								},
							)
						}
					}
				}(secret, mtMap)
			}
		}
		// 每个频率区间遍历一次
		time.Sleep(time.Duration(self.MessageDuration) * time.Second)
	}
}

// 消息合并， 若消息能合并，则rim2 = nil; 若消息不能合并rim1, 为原始消息， rim2为其余消息
// rim1: 合并后的消息
// rim2: 不能合并的消息列表
func (self *BotPool) mergeMessage(im1, im2 IMessage) (rim1 IMessage, rim2 []IMessage) {
	rim := im1.MergeMessage(self.Boundary, im1, im2)
	if len(rim) == 1 {
		rim1 = rim[0]
		return
	}

	if len(rim) > 1 {
		rim1 = rim[0]
		rim2 = rim[1:]
		return
	}

	return
}
