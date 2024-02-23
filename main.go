package main

import (
	"context"
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/scroll-tech/go-ethereum/ethclient"
	"github.com/scroll-tech/go-ethereum/log"
	"github.com/scroll-tech/go-ethereum/rpc"
)

var l2GethEndpoint = "http://localhost:8545"

func init() {
	output := io.Writer(os.Stderr)
	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	ostream := log.StreamHandler(output, log.TerminalFormat(usecolor))
	glogger := log.NewGlogHandler(ostream)
	// Set log level
	glogger.Verbosity(log.LvlInfo)
	log.Root().SetHandler(glogger)
}

func main() {
	ctx := context.Background()

	l2GethClient, err := ethclient.DialContext(ctx, l2GethEndpoint)
	if err != nil {
		panic(err)
	}
	// Use gzip compression.
	l2GethClient.SetHeader("Accept-Encoding", "gzip")

	// GetNumSkippedTransactions
	nSkipped, err := l2GethClient.GetNumSkippedTransactions(ctx)
	if err != nil {
		panic(err)
	}
	log.Info("GetNumSkippedTransactions", "nSkipped", nSkipped)

	// GetSkippedTransactionHashes
	hashList, err := l2GethClient.GetSkippedTransactionHashes(ctx, 0, uint64(nSkipped))
	if err != nil {
		panic(err)
	}

	for i, hash := range hashList {
		log.Info("handling tx", "hash", hash, "i", i, "total", nSkipped)

		// GetSkippedTransaction
		tx, err := l2GethClient.GetSkippedTransaction(ctx, hash)
		if err != nil {
			panic(err)
		}

		// GetTxBlockTraceOnTopOfBlock
		blockNum := rpc.BlockNumber(0)
		_, err = l2GethClient.GetTxBlockTraceOnTopOfBlock(ctx, tx, rpc.BlockNumberOrHash{BlockNumber: &blockNum} /*TODO*/, nil)

		// dump trace
	}
}
