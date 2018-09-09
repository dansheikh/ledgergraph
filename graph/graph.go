package graph

type Transaction struct {
	Vendor   string  `json:"vendor"`
	Customer string  `json:"customer"`
	Amount   float64 `json:"amount"`
}

type Data struct {
	Index        int           `json:"index"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PreviousHash string        `json:"previousHash"`
	Nonce        int           `json:"nonce"`
}

type Metadata struct {
	Hash     string  `json:"hash"`
	Previous *Vertex `json:"previous"`
	Next     *Vertex `json:"next"`
}

type Vertex struct {
	Data     Data     `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type AcyclicGraph struct {
	Head *Vertex `json:"head"`
	Tail *Vertex `json:"tail"`
}
