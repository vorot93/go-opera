package launcher

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Fantom-foundation/lachesis-base/hash"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/kvdb/batched"
	"github.com/Fantom-foundation/lachesis-base/kvdb/pebble"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/utils/dbutil/autocompact"
)

var (
	eventsFileHeader  = hexutils.HexToBytes("7e995678")
	eventsFileVersion = hexutils.HexToBytes("00010001")
)

// statsReportLimit is the time limit during import and export after which we
// always print out progress. This avoids the user wondering what's going on.
const statsReportLimit = 8 * time.Second

func exportEvents(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	cfg := makeAllConfigs(ctx)

	rawDbs := makeDirectDBsProducer(cfg)
	gdb := makeGossipStore(rawDbs, cfg)
	defer gdb.Close()

	fn := ctx.Args().First()

	// Open the file handle and potentially wrap with a gzip stream
	fh, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fh.Close()

	var writer io.Writer = fh
	if strings.HasSuffix(fn, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	from := idx.Epoch(1)
	if len(ctx.Args()) > 1 {
		n, err := strconv.ParseUint(ctx.Args().Get(1), 10, 32)
		if err != nil {
			return err
		}
		from = idx.Epoch(n)
	}
	to := idx.Epoch(0)
	if len(ctx.Args()) > 2 {
		n, err := strconv.ParseUint(ctx.Args().Get(2), 10, 32)
		if err != nil {
			return err
		}
		to = idx.Epoch(n)
	}

	log.Info("Exporting events to file", "file", fn)
	// Write header and version
	_, err = writer.Write(append(eventsFileHeader, eventsFileVersion...))
	if err != nil {
		return err
	}
	err = exportTo(writer, gdb, from, to)
	if err != nil {
		utils.Fatalf("Export error: %v\n", err)
	}

	return nil
}

// exportTo writer the active chain.
func exportTo(w io.Writer, gdb *gossip.Store, from, to idx.Epoch) (err error) {
	start, reported := time.Now(), time.Time{}

	var (
		counter int
		last    hash.Event
	)
	gdb.ForEachEventRLP(from.Bytes(), func(id hash.Event, event rlp.RawValue) bool {
		if to >= from && id.Epoch() > to {
			return false
		}
		counter++
		_, err = w.Write(event)
		if err != nil {
			return false
		}
		last = id
		if counter%100 == 1 && time.Since(reported) >= statsReportLimit {
			log.Info("Exporting events", "last", last.String(), "exported", counter, "elapsed", common.PrettyDuration(time.Since(start)))
			reported = time.Now()
		}
		return true
	})
	log.Info("Exported events", "last", last.String(), "exported", counter, "elapsed", common.PrettyDuration(time.Since(start)))

	return
}

func exportEvmKeys(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	cfg := makeAllConfigs(ctx)

	rawDbs := makeDirectDBsProducer(cfg)
	gdb := makeGossipStore(rawDbs, cfg)
	defer gdb.Close()

	fn := ctx.Args().First()

	keysDB_, err := pebble.New(fn, 1024*opt.MiB, utils.MakeDatabaseHandles()/2, nil, nil)
	if err != nil {
		return err
	}
	keysDB := batched.Wrap(autocompact.Wrap2M(keysDB_, opt.GiB, 16*opt.GiB, true, "evm-keys"))
	defer keysDB.Close()

	it := gdb.EvmStore().EvmDb.NewIterator(nil, nil)
	// iterate only over MPT data
	it = mptAndPreimageIterator{it}
	defer it.Release()

	log.Info("Exporting EVM keys", "dir", fn)
	for it.Next() {
		if err := keysDB.Put(it.Key(), []byte{0}); err != nil {
			return err
		}
	}
	log.Info("Exported EVM keys", "dir", fn)
	return nil
}

func exportAtroposes(ctx *cli.Context) error {
	cfg := makeAllConfigs(ctx)

	rawDbs := makeDirectDBsProducer(cfg)
	gdb := makeGossipStore(rawDbs, cfg)
	defer gdb.Close()

	fileName := ctx.String(CSVFileFlag.Name)

	log.Info("Exporting atroposes data into", "file", fileName)
	const batchSize = 1_000_000
	i := 0

	// setup writer
	csvOut, err := os.Create(fileName)
	if err != nil {
		panic(fmt.Sprintf("Unable to create file %s.", fileName))
	}
	defer csvOut.Close()

	w := csv.NewWriter(csvOut)

	var headers []string
	if ctx.GlobalBool(CSVExtendedFlag.Name) {
		headers = []string{
			"epoch",
			"blockNumber",
			"atropos",
			"txs",
			"gasused",
		}
	} else {
		headers = []string{
			"blockNumber",
			"atropos",
		}
	}
	if err = w.Write(headers); err != nil {
		panic(err)
	}

	buffered := 0
	var latest uint64

	minEpoch := ctx.Uint64(MinEpochFlag.Name)
	maxEpoch := ctx.Uint64(MaxEpochFlag.Name)

	gdb.ForEachBlock(func(index idx.Block, block *inter.Block) {
		epoch := block.Atropos.Epoch()
		if uint64(epoch) >= minEpoch && uint64(epoch) <= maxEpoch {
			var line []string
			if ctx.GlobalBool(CSVExtendedFlag.Name) {
				line = []string{
					strconv.Itoa(int(epoch)),
					strconv.Itoa(int(index)),
					block.Atropos.Hex(),
					strconv.Itoa(len(block.InternalTxs) + len(block.SkippedTxs) + len(block.Txs)),
					strconv.Itoa(int(block.GasUsed)),
				}
			} else {
				line = []string{
					strconv.Itoa(int(index)),
					block.Atropos.Hex(),
				}
			}

			err = w.Write(line)
			if err != nil {
				panic(fmt.Sprintf("block %d: %v", index, err))
			}

			i++
			buffered++
			latest = uint64(index)

			if i%batchSize == 0 {
				w.Flush()
				buffered = 0

				log.Info("Exporting atroposes data into", "file", fileName, "block", index)
			}
		}
	})

	if buffered != 0 {
		w.Flush()
		log.Info("Exporting atroposes data into", "file", fileName, "block", latest)
	}

	log.Info("Exported atroposes data into", "file", fileName)

	return nil
}
