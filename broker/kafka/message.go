package kafka

import (
	kafkago "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*MessageCarrier)(nil)

type MessageCarrier struct {
	msg *kafkago.Message
}

func NewMessageCarrier(msg *kafkago.Message) MessageCarrier {
	return MessageCarrier{msg: msg}
}

func (c MessageCarrier) Get(key string) string {
	for _, h := range c.msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c MessageCarrier) Set(key, val string) {
	for i := 0; i < len(c.msg.Headers); i++ {
		if c.msg.Headers[i].Key == key {
			c.msg.Headers = append(c.msg.Headers[:i], c.msg.Headers[i+1:]...)
			i--
		}
	}
	c.msg.Headers = append(c.msg.Headers, kafkago.Header{
		Key:   key,
		Value: []byte(val),
	})
}

func (c MessageCarrier) Keys() []string {
	out := make([]string, len(c.msg.Headers))
	for i, h := range c.msg.Headers {
		out[i] = h.Key
	}
	return out
}
