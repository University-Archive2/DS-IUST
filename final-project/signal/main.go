package signal

func main() {
	kafkaProducer := NewKafkaProducer("stock", configs.KafkaHosts, 2)
}
