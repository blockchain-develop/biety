package main

import (
	"fmt"
	"github.com/biety/blocksync"
	"github.com/biety/consensus"
	"github.com/biety/jsonrpc"
	"github.com/biety/ledger"
	"github.com/biety/p2pserver"
	"github.com/biety/txnpool"
	"github.com/biety/validator"
	"github.com/ontio/ontology-eventbus/actor"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "biety cli"
	app.Action = startBiety
	app.Version = "0.1"
	app.Copyright = "Copyright in 2018 The biety Authors"
	app.Commands = []cli.Command {

	}
	app.Flags = []cli.Flag {

	}

	return app
}

func startBiety(ctx* cli.Context) {

	fmt.Printf("init Ledger\n")
	ldg, err := initLedger(ctx)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer ldg.Close()

	fmt.Printf("init TxPool\n")
	txnpoolserver, err := initTxPool(ctx)
	if err != nil {
		fmt.Print(err)
		return
	}


	fmt.Printf("start p2p networks\n")
	p2p, p2pactorpid, err := initP2PNode(ctx)
	//p2p, _, err := initP2PNode(ctx)
	if err != nil {
		fmt.Print(err)
		return
	}
	p2pserver.SetTxnPoolPid(txnpoolserver.GetTxActorPID())

	fmt.Printf("init block\n")
	_, err = initBlockSync(ctx, p2p, ldg)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("init consensus\n")
	cons, err := initConsensus(ctx)
	//_, err = initConsensus(ctx)
	if err != nil {
		fmt.Print(err)
		return
	}
	p2pserver.SetConsensusPid(cons.GetPID())
	cons.Init(p2pactorpid, txnpoolserver.GetTxPoolActorPID())
	cons.Start()


	fmt.Printf("init rpc\n")
	err = initRpc(ctx)
	if err != nil {
		fmt.Print(err)
		return
	}


	waitToExit()
}

func initLedger(ctx *cli.Context) (*ledger.Ledger, error) {
	ldg, err := ledger.NewLedger("")
	return ldg, err
}

func initTxPool(ctx *cli.Context) (*txnpool.TxPoolServer, error) {
	txnpoolserver, err := txnpool.StartTxnPoolServer()
	if err != nil {
		return nil, err
	}

	v1, _ := validator.NewValidator("stateless_validator")
	v1.Register(txnpoolserver.GetVerifyRspActorPID())

	v2, _ := validator.NewValidator("stateless_validator2")
	v2.Register(txnpoolserver.GetVerifyRspActorPID())

	v3, _ := validator.NewValidator("stateful_validator")
	v3.Register(txnpoolserver.GetVerifyRspActorPID())
	return txnpoolserver, nil
}

func initP2PNode(ctx *cli.Context) (*p2pserver.P2PServer, *actor.PID, error) {
	p2p := p2pserver.NewServer()
	p2pactor := p2pserver.NewP2PActor()
	p2pactorid, err := p2pactor.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("p2pActor init error %s", err)
	}
	p2p.SetPID(p2pactorid)
	err = p2p.Start()
	if err != nil {
		return nil,nil, fmt.Errorf("init P2P failed, err %s", err)
	}

	return p2p, p2pactorid, nil
}

func initBlockSync(ctx *cli.Context, server *p2pserver.P2PServer, ledger *ledger.Ledger) (*blocksync.BlockSyncMgr, error) {
	blocksync := blocksync.NewBlockSyncMgr(server, ledger)
	go blocksync.Start()
	return blocksync, nil
}

func initConsensus(ctx *cli.Context) (consensus.ConsensusService, error) {
	s, err := consensus.NewConsensueService()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func initRpc(ctx *cli.Context) error {
	err := jsonrpc.StartRPCServer()
	if err != nil {
		return err
	}
	return nil
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			fmt.Printf("biety received exit signal : %v", sig)
			close(exit)
			break
		}
	}()

	<-exit
}
