package models

// Conta representa a estrutura dos dados de uma conta para o JSON e para o banco de dados
type Conta struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
	Industry string `json:"industry"`
}

// Resposta é a estrutura padrão para as respostas da API
type Resposta struct {
	Status   string `json:"status"`
	Mensagem string `json:"mensagem"`
}
