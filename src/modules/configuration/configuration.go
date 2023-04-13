package configuration

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config") // nome do arquivo que queremos carregar
	viper.SetConfigType("env")    // extensão do arquivo
	viper.AddConfigPath("./conf") // caminho alternativo onde está o arquivo
	err := viper.ReadInConfig()   // lê o arquivo e carrega seu conteúdo
	if err != nil {
		panic(err)
	}
}

func Get(key string) string {
	return viper.GetString("port")
}
