package ledger

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/dansheikh/ledgergraph/graph"
)

/*
CreateTransaction creates a new transaction pending inclusion into the next vertex.
*/
func CreateTransaction(transactions []graph.Transaction, vendor string, customer string, amount float64) []graph.Transaction {
	newTransaction := graph.Transaction{Vendor: vendor, Customer: customer, Amount: amount}
	transactions = append(transactions, newTransaction)
	return transactions
}

/*
CreateVertex creates a new vertex and updates directed acyclic graph accordingly.
*/
func CreateVertex(transactions []graph.Transaction, acyclicGraph *graph.AcyclicGraph) string {
	index := 1
	previousHash := ""
	head := (*acyclicGraph).Head
	tail := (*acyclicGraph).Tail

	if tail != nil {
		index = (*tail).Data.Index + 1
		previousHash = (*tail).Metadata.Hash
	}

	vertexData := graph.Data{Index: index, Timestamp: time.Now().Format(time.RFC3339), Transactions: transactions, PreviousHash: previousHash}
	vertexMetadata := graph.Metadata{}

	newVertex := graph.Vertex{Data: vertexData, Metadata: vertexMetadata}

	// Determine nonce for appropriate hash.
	hash := pow(&newVertex)
	newVertex.Metadata.Hash = hash

	if head == nil && tail == nil {
		(*acyclicGraph).Head = &newVertex
		(*acyclicGraph).Tail = &newVertex
	} else {
		// Add edge from/to previous vertex.
		(*tail).Metadata.Next = &newVertex
		newVertex.Metadata.Previous = tail

		// Update tail.
		(*acyclicGraph).Tail = &newVertex
	}

	return newVertex.Metadata.Hash
}

func pow(vertex *graph.Vertex) string {
	var shaSum string

	for {
		vertexString := fmt.Sprintf("%v", (*vertex).Data)
		tmp := sha256.Sum256([]byte(vertexString))
		shaSum = base64.URLEncoding.EncodeToString(tmp[:])

		if strings.HasPrefix(shaSum, "0000") {
			break
		}

		(*vertex).Data.Nonce = (*vertex).Data.Nonce + 1
	}

	return shaSum
}

/*
ValidateLedger validates ledger graph.
*/
func ValidateLedger(graph *graph.AcyclicGraph) (bool, int) {
	vertex := (*graph).Head
	valid := true

	for {
		nonce := (*vertex).Data.Nonce
		hash := pow(vertex)
		next := (*vertex).Metadata.Next

		if nonce != (*vertex).Data.Nonce || hash != (*vertex).Metadata.Hash {
			valid = false
			break
		}
		if next == nil {
			break
		} else if hash != (*next).Data.PreviousHash {
			valid = false
			break
		} else {
			vertex = next
		}
	}

	return valid, (*vertex).Data.Index
}
