package config

import (
	"path/filepath"
)

// Definimos la ruta de configuración base para Windows
const baseConfigPath = "C:/Users/danie/Documents/UP Daniel/Computo Distribuido/GO_Server/GO_Module/0250952_SistemasDistribuidos/test"

var (
	CAFile               = configFile("ca.pem")
	ServerCertFile       = configFile("server.pem")
	ServerKeyFile        = configFile("server-key.pem")
	RootClientCertFile   = configFile("root-client.pem")
	RootClientKeyFile    = configFile("root-client-key.pem")
	NobodyClientCertFile = configFile("nobody-client.pem")
	NobodyClientKeyFile  = configFile("nobody-client-key.pem")
	ACLModelFile         = configFile("model.conf")
	ACLPolicyFile        = configFile("policy.csv")
)

func configFile(filename string) string {
	// Devuelve la ruta del archivo en el directorio específico de Windows
	return filepath.Join(baseConfigPath, filename)
}
