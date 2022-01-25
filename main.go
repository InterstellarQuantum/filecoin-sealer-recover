package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/filecoin-project/go-address"
	cliutil "github.com/filecoin-project/lotus/cli/util"
	"github.com/froghub-io/filecoin-sealer-recover/recovery"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"os"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	app := &cli.App{
		Name:    "sealer-recovery",
		Usage:   "Filecoin sealer recovery",
		Version: BuildVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "miner",
				Usage:    "Filecoin miner. Such as: f01000",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "cid",
				Usage:    "PreCommitCID~~",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "ticket",
				Usage:    "ticket~~",
				Required: true,
			},
			&cli.IntSliceFlag{
				Name:     "sector",
				Usage:    "Sector number to be recovered. Such as: --sector=0 --sector=1 ... ",
				Required: true,
			},
			&cli.UintFlag{
				Name:  "parallel",
				Usage: "Number of parallel P1",
				Value: 1,
			},
			&cli.StringFlag{
				Name:  "sealing-result",
				Value: "/lotus_data/store",
				Usage: "Recover sector result path",
			},
			&cli.StringFlag{
				Name:  "sealing-temp",
				Value: "/lotus_data/temp",
				Usage: "Temporarily generated during sector recovery",
			},
			&cli.Int64Flag{
				Name:  "rsp",
				Value: 8,
				Usage: "扇区版本号默认是v1_1版本32G,v1_1版本64G是9,v1版本32G是3,64G是4,如果不是v1_1版本32G扇区,需要手动设置该值",
			},
			&cli.Uint64Flag{
				Name:  "size",
				Value: 32,
				Usage: "扇区大小,默认32",
			},
		},
		Action: func(cctx *cli.Context) error {
			log.Info("Start sealer recovery!")
			ctx := cliutil.DaemonContext(cctx)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			recovery.Ss = cctx.String("cid") //cid的string
			tip := cctx.String("ticket")
			bytes, e := hex.DecodeString(tip) //扇区的ticket
			if e != nil {
				fmt.Printf("输入的ticket错误, %s", e)
				return nil
			}
			recovery.Rns = bytes
			maddr, err := address.NewFromString(cctx.String("miner"))
			if err != nil {
				return xerrors.Errorf("Getting NewFromString err:", err)
			}
			actorID, err := address.IDFromAddress(maddr)
			if err != nil {
				return xerrors.Errorf("Getting IDFromAddress err:", err)
			}
			if err = recovery.RecoverSealedFile(ctx, maddr, actorID, cctx.IntSlice("sector"), cctx.Uint("parallel"), cctx.String("sealing-result"), cctx.String("sealing-temp"), cctx.Int64("rsp"), cctx.Uint64("size")); err != nil {
				return err
			}
			log.Info("Complete recovery sealed!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Warnf("%+v", err)
		return
	}

}

const BuildVersion = "1.1.0"
