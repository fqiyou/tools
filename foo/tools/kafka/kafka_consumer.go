package kafka

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/fqiyou/tools/foo/tools/logs"
	"strings"
)


type KafkaConsumer struct {
	consumer *cluster.Consumer
	stopped  chan struct{}
	msgs     chan ([]byte)

	Name          string
	Brokers       string
	ConsumerGroup string
	Topic         string

	Sasl struct {
		Username string
		Password string
	}
}

func NewKafkaConsumer() *KafkaConsumer {
	return &KafkaConsumer{}
}

func (k *KafkaConsumer) Init() error {
	k.msgs = make(chan []byte, 300000)
	k.stopped = make(chan struct{})
	return nil
}

func (k *KafkaConsumer) Msgs() chan []byte {
	return k.msgs
}

func (k *KafkaConsumer) Start() error {

	config := cluster.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	//config.Consumer.Offsets.Initial = sarama.OffsetOldest
	if k.Sasl.Username != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = k.Sasl.Username
		config.Net.SASL.Password = k.Sasl.Password
	}
	consumer, err := cluster.NewConsumer(strings.Split(k.Brokers, ","), k.ConsumerGroup, []string{k.Topic}, config)
	if err != nil {
		return err
	}
	k.consumer = consumer
	go func() {
		log.Info("Start kafka services", k.Name)
	FOR:
		for {
			select {
			case msg, more := <-k.consumer.Messages():
				if !more {
					break FOR
				}
				k.msgs <- msg.Value
				k.consumer.MarkOffset(msg, "") // mark message as processed

			case err, more := <-k.consumer.Errors():
				if more {
					log.Error("Error: %s\n", err.Error())
				}
			}
		}
		close(k.stopped)
	}()
	return nil
}

func (k *KafkaConsumer) Stop() error {
	k.consumer.Close()
	<-k.stopped
	close(k.msgs)
	return nil
}

func (k *KafkaConsumer) Description() string {
	return "kafka consumer:" + k.Topic
}

func (k *KafkaConsumer) GetName() string {
	return k.Name
}
