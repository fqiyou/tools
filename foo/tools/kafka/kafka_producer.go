package kafka

import (
	"github.com/Shopify/sarama"
)

//import "github.com/bsm/sarama-cluster"

type KafkaProducer struct {
	producer *sarama.SyncProducer

	Name          string
	Brokers       string
	Topic         string

	Sasl struct {
		Username string
		Password string
	}
}


func NewKafkaProducer() *KafkaConsumer {
	return &KafkaConsumer{}
}

func (k *KafkaProducer) Init() error {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	if k.Sasl.Username != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = k.Sasl.Username
		config.Net.SASL.Password = k.Sasl.Password
	}
	//
	//k.producer ,err= sarama.NewSyncProducer(strings.Split(k.Brokers, ","), config)
	//if err != nil {
	//	return err
	//}
	return nil
}


//
//func (k *KafkaConsumer) Stop() error {
//	k.consumer.Close()
//	<-k.stopped
//	close(k.msgs)
//	return nil
//}
//
