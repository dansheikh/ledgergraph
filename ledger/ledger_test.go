package ledger

import (
	"encoding/json"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/dansheikh/ledgergraph/graph"
)

func graphFixture() *graph.AcyclicGraph {
	acyclicGraph := graph.AcyclicGraph{}
	init := make([]graph.Transaction, 0)

	for i := 0; i < 3; i++ {
		num := math.Round(rand.Float64()*10000000) / 100
		tmp := CreateTransaction(init, "Pepsi", "Taco Bell", num)
		CreateVertex(tmp, &acyclicGraph)
	}

	return &acyclicGraph
}

func TestCreateTransaction(t *testing.T) {
	init := make([]graph.Transaction, 0)
	tmp := CreateTransaction(init, "Pepsi", "Taco Bell", 86931.27)
	transactionSize := len(tmp)

	if transactionSize != 1 {
		t.Error("Expected transaction slice of 1, received", transactionSize)
	}

	transaction := tmp[0]

	if transaction.Vendor != "Pepsi" || transaction.Customer != "Taco Bell" || transaction.Amount != 86931.27 {
		transactionJSON, _ := json.Marshal(transaction)
		t.Error("Expected transaction with {\"vendor\": \"Pepsi\", \"customer\": \"Taco Bell\", \"amount\": 86931.21}, received", string(transactionJSON))
	}
}
func TestCreateVertex(t *testing.T) {
	acyclicGraph := graph.AcyclicGraph{}
	init := make([]graph.Transaction, 0)
	tmp := CreateTransaction(init, "Pepsi", "Taco Bell", 86931.27)

	hash := CreateVertex(tmp, &acyclicGraph)

	if len(hash) == 0 {
		t.Error("Expected non-zero hash, received hash of length", len(hash))
	}

	if !strings.HasPrefix(hash, "0000") {
		t.Error("Expected (hash) string prefixed with \"0000\", received", hash[0:2])
	}
}

func TestValidateGraph(t *testing.T) {
	acyclicGraph := graphFixture()

	validity, count := ValidateLedger(acyclicGraph)

	if validity != true {
		t.Error("Expected acyclic graph to be (valid) true, received (invalid)", validity)
	}

	if count != 3 {
		t.Error("Expected acyclic graph with 3 vertices, received graph with", count, "vertices")
	}
}
