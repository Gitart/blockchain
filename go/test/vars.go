package main

const(
     GlobalData = "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
)

type Mst map[string] interface{}

var (
     GlobalPass = ""
)


type Go struct {
	Jsonrpc string `json:"jsonrpc" db:"jsonrpc"`
	ID int `json:"id" db:"id"`
	Result struct {
		Difficulty           string `json:"difficulty" db:"difficulty"`
		ExtraData            string `json:"extraData" db:"extra_data"`
		GasLimit             string `json:"gasLimit" db:"gas_limit"`
		GasUsed              string `json:"gasUsed" db:"gas_used"`
		Hash                 string `json:"hash" db:"hash"`
		LogsBloom            string `json:"logsBloom" db:"logs_bloom"`
		Miner                string `json:"miner" db:"miner"`
		MixHash              string `json:"mixHash" db:"mix_hash"`
		Nonce                string `json:"nonce" db:"nonce"`
		Number               string `json:"number" db:"number"`
		ParentHash           string `json:"parentHash" db:"parent_hash"`
		ReceiptsRoot         string `json:"receiptsRoot" db:"receipts_root"`
		Sha3Uncles           string `json:"sha3Uncles" db:"sha3_uncles"`
		Size                 string `json:"size" db:"size"`
		StateRoot            string `json:"stateRoot" db:"state_root"`
		Timestamp            string `json:"timestamp" db:"timestamp"`
		TotalDifficulty      string `json:"totalDifficulty" db:"total_difficulty"`
		Transactions []struct {
			BlockHash        string `json:"blockHash" db:"block_hash"`
			BlockNumber      string `json:"blockNumber" db:"block_number"`
			From             string `json:"from" db:"from"`
			Gas              string `json:"gas" db:"gas"`
			GasPrice         string `json:"gasPrice" db:"gas_price"`
			Hash             string `json:"hash" db:"hash"`
			Input            string `json:"input" db:"input"`
			Nonce            string `json:"nonce" db:"nonce"`
			To               string `json:"to" db:"to"`
			TransactionIndex string `json:"transactionIndex" db:"transaction_index"`
			Value            string `json:"value" db:"value"`
			V                string `json:"v" db:"v"`
			R                string `json:"r" db:"r"`
			S                string `json:"s" db:"s"`
		} `json:"transactions" db:"transactions"`
		TransactionsRoot     string `json:"transactionsRoot" db:"transactions_root"`
		Uncles               []interface{} `json:"uncles" db:"uncles"`} `json:"result" db:"result"`
}
