package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/fqiyou/tools/foo/tools/message"
	"github.com/fqiyou/tools/foo/util"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"

	//"sync"
	"testing"
)

func TestNewKafkaConsumer(t *testing.T) {
	fmt.Println("开始")
	os.Exit(0)
	util.Log.SetLevel(logrus.ErrorLevel)

	kafka_consumer := NewKafkaConsumer()
	kafka_consumer.Brokers = "spark016:9092"
	//kafka_consumer.Brokers = "bj-dcs-005:9092"
	kafka_consumer.ConsumerGroup = "dev_yc_test_0129"
	kafka_consumer.Topic = "locker" // article
	kafka_consumer.Name = "test-yc-article"

	kafka_consumer.Init()
	kafka_consumer.Start()

	defer kafka_consumer.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	message := message.NewMessage()




	topic_name := "yc_test"
	consumer_string := "spark016:9092,spark017:9092,spark018:9092"

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_0
	//config.Producer.MaxMessageBytes = 104857500
	producer, err := sarama.NewAsyncProducer(strings.Split(consumer_string, ","), config)


	if err != nil {
		util.Log.Error(err)
	}

	util.Log.Info("start make producer!!!")


	defer producer.AsyncClose()
	go func(p sarama.AsyncProducer) {
		for{
			select {
			//case  suc := <-p.Successes():
				//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
			case fail := <-p.Errors():
				util.Log.Error("err: ", fail.Err)

			}
		}
	}(producer)

	//producer_msg_chan := make(chan []byte, 300000)

	//var wg    sync.WaitGroup


FOR:
	for {
		select {
		case msg, more := <-kafka_consumer.Msgs():
			if !more {
				break FOR
			}

			go func() {

				//t1 := time.Now()
				msg_list := message.ToMessageList(string(msg))

				//elapsed := time.Since(t1)
				//util.Log.Info(elapsed)

				for _, v := range msg_list {
					//log.Error(v)
					//util.JsonPrint(util.ModelToString(v))
					//util.JsonPrint(v)
					producer_msg := &sarama.ProducerMessage{
						Topic: topic_name,
					}

					producer_msg.Value = sarama.ByteEncoder(util.ModelToString(v))

					//util.Log.Info(util.ModelToString(v))

					//wg.Add(1)
					//producer_msg_chan <- sarama.ByteEncoder(util.ModelToString(v))
					producer.Input() <- producer_msg
				}
				//msg_list = nil
				//select {
				//case suc := <-producer.Successes():
				//	log.Info("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
				//case fail := <-producer.Errors():
				//	log.Error("err: ", fail.Err)
				//}

			}()
			//case <-c:
			//	kafka_consumer.Stop()

		}

	}
	//producer_msg = nil
	//
	//
	//for v := range producer_msg_chan {
	//	//producer_msg.Value = sarama.ByteEncoder(util.ModelToString(v))
	//	//producer.Input() <- producer_msg
	//	select {
	//	case suc := <-producer.Successes():
	//		util.Log.Info("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
	//	case fail := <-producer.Errors():
	//		util.Log.Error("err: ", fail.Err)
	//	}
	//}



}
