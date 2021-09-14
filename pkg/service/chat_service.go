package service

//Objeto que contiene a la interfaz propia de la base de datos
type serviceProperties struct {
	repository ChatRepository
}

type ChatService interface {
}

//Inicializa el servicio que se comunica con la base de datos
func NewChatService(repository ChatRepository) ChatService {
	return &serviceProperties{repository}
}
