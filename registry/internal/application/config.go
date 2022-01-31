package application

type Config struct {
	ApplicationHost    string `envconfig:"HOST" default:"localhost" required:"true"`
	ApplicationPort    int    `envconfig:"PORT" default:"80" required:"true"`
	DiagnosticPort     int    `envconfig:"DIAG_PORT" default:"81" required:"true"`
	MongoConnectionUrl string `envconfig:"MONGO_URL" default:"mongodb://root:password@mongodb-primary:27017/" required:"true"`
	KafkaConnectionUrl string `envconfig:"KAFKA_URL" default:"kafka:9092" required:"true"`

	//MongoConnectionUrl string `envconfig:"MONGO_URL" default:"mongodb://root:password@localhost:27017/" required:"true"`
	//KafkaConnectionUrl string `envconfig:"KAFKA_URL" default:"localhost:9093" required:"true"`
}
