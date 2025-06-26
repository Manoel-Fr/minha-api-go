package models

// Conta representa a estrutura dos dados de uma conta para o JSON e para o banco de dados
type Conta struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
	Cpf      string `json:"cpf"`
}

// Resposta é a estrutura padrão para as respostas da API
type Resposta struct {
	Status   string `json:"status"`
	Mensagem string `json:"mensagem"`
}
