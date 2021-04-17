package locale

import (
	"errors"
	"os"
	"sync"
)

// MsgSet include of  a message sets
type MsgSet map[string]string

type lang string

const (
	// ENUS for English in USA style
	ENUS lang = "en_US.UTF-8"
	// ZHCN for Chinese in China style
	ZHCN lang = "zh_CN.UTF-8"
)

// Localer include of system languange and message set
// In linux the value of lang depend on systemctl env 'LANG'
type Localer struct {
	mu      sync.RWMutex
	lang    lang
	msgSets map[lang]MsgSet
}

// GetLang get language of system
// return a string like locale
// ex.en_US.utf8
// zh_CN.utf8
func (l *Localer) GetLang() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return string(l.lang)
}

//SetMsgs set messages for tagert language
func (l *Localer) SetMsgs(lan lang, mS MsgSet) error {
	if len(mS) < 1 {
		return errors.New("MsgSet is empty")
	}
	if _, exist := l.msgSets[lan]; exist {
		return errors.New("Message set existed in this language:" + string(lan))
	}
	msgSet := make(MsgSet)
	for k, v := range mS {
		msgSet[k] = v
	}
	l.msgSets[lan] = msgSet
	return nil
}

// AppendMsg append a message set in prticular language
func (l *Localer) AppendMsg(lan lang, mS MsgSet) error {
	if len(mS) < 1 {
		return errors.New("MsgSet is empty")
	}
	if _, exist := l.msgSets[lan]; !exist {
		return errors.New("Message set isn't existed in this language:" + string(lan))
	}
	for k, v := range mS {
		l.msgSets[lan][k] = v
	}
	return nil
}

// GetMsgWithError return a message in system language kind
// If message sets isn't exist it will return message in English format
// return error if msgKey isn't exist in msgSet
func (l *Localer) GetMsgWithError(msgKey string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	lang := l.lang
	if _, langExist := l.msgSets[l.lang]; !langExist {
		lang = ENUS
	}
	if _, msgExist := l.msgSets[lang][msgKey]; !msgExist {
		return "", errors.New("message key isn't exist")
	}
	return l.msgSets[lang][msgKey], nil
}

//GetMsg return message in system lang kind
//return english kind of message if system language isn.t exist in msgSets
//If messages ket isn't exist it will return ""
func (l *Localer) GetMsg(msgKey string) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	lang := l.lang
	if _, langExist := l.msgSets[l.lang]; !langExist {
		lang = ENUS
	}
	if _, msgExist := l.msgSets[lang][msgKey]; !msgExist {
		return ""
	}
	return l.msgSets[lang][msgKey]
}

// New implement a Localer
func New() *Localer {
	sysLang := os.Getenv("LANG")
	return &Localer{
		lang:    lang(sysLang),
		msgSets: make(map[lang]MsgSet),
	}
}
