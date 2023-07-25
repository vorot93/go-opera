package ethapi

import (
	"context"

	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// PublicAbftAPI provides an API to access consensus related information.
// It offers only methods that operate on public data that is freely available to anyone.
type PublicAbftAPI struct {
	b Backend
}

// NewPublicAbftAPI creates a new SFC protocol API.
func NewPublicAbftAPI(b Backend) *PublicAbftAPI {
	return &PublicAbftAPI{b}
}

type EventInfo struct {
	ID           common.Hash    `json:"id"`
	GasPowerLeft GasPowerLeft   `json:"gasPowerLeft"`
	Time         hexutil.Uint64 `json:"time"`
}

type ValidatorEpochState struct {
	GasRefund      hexutil.Uint64 `json:"gasRefund"`
	PrevEpochEvent EventInfo      `json:"prevEpochEvent"`
}

type GasPowerLeft [2]hexutil.Uint64

func (s *PublicAbftAPI) GetFullEpochState(ctx context.Context, epoch rpc.BlockNumber) (map[string]interface{}, error) {
	_, es, err := s.b.GetEpochBlockState(ctx, epoch)
	if err != nil {
		return nil, err
	}
	if es == nil {
		return nil, nil
	}

	validators := map[hexutil.Uint64]interface{}{}
	profiles := es.ValidatorProfiles
	for _, vid := range es.Validators.IDs() {
		validators[hexutil.Uint64(vid)] = map[string]interface{}{
			"weight": (*hexutil.Big)(profiles[vid].Weight),
			"pubkey": profiles[vid].PubKey.String(),
		}
	}

	validator_state := make([]ValidatorEpochState, len(es.ValidatorStates))
	for i, vs := range es.ValidatorStates {
		gas_power_left := GasPowerLeft{}
		for j, gas := range vs.PrevEpochEvent.GasPowerLeft.Gas {
			gas_power_left[j] = hexutil.Uint64(gas)
		}
		validator_state[i] = ValidatorEpochState{
			GasRefund: hexutil.Uint64(vs.GasRefund),
			PrevEpochEvent: EventInfo{
				ID:           common.Hash(vs.PrevEpochEvent.ID),
				GasPowerLeft: gas_power_left,
				Time:         hexutil.Uint64(vs.PrevEpochEvent.Time),
			},
		}
	}

	res := map[string]interface{}{
		"epoch":              hexutil.Uint64(epoch),
		"epoch_start":        hexutil.Uint64(es.EpochStart),
		"prev_epoch_start":   hexutil.Uint64(es.PrevEpochStart),
		"epoch_state_root":   es.EpochStateRoot.String(),
		"validators":         nil,
		"validator_profiles": validators,
		"validator_state":    validator_state,
	}

	return res, nil

}

func (s *PublicAbftAPI) GetValidators(ctx context.Context, epoch rpc.BlockNumber) (map[hexutil.Uint64]interface{}, error) {
	bs, es, err := s.b.GetEpochBlockState(ctx, epoch)
	if err != nil {
		return nil, err
	}
	if es == nil {
		return nil, nil
	}
	res := map[hexutil.Uint64]interface{}{}
	for _, vid := range es.Validators.IDs() {
		profiles := es.ValidatorProfiles
		if epoch == rpc.PendingBlockNumber {
			profiles = bs.NextValidatorProfiles
		}
		res[hexutil.Uint64(vid)] = map[string]interface{}{
			"weight": (*hexutil.Big)(profiles[vid].Weight),
			"pubkey": profiles[vid].PubKey.String(),
		}
	}
	return res, nil
}

// GetDowntime returns validator's downtime.
func (s *PublicAbftAPI) GetDowntime(ctx context.Context, validatorID hexutil.Uint) (map[string]interface{}, error) {
	blocks, period, err := s.b.GetDowntime(ctx, idx.ValidatorID(validatorID))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"offlineBlocks": hexutil.Uint64(blocks),
		"offlineTime":   hexutil.Uint64(period),
	}, nil
}

// GetEpochUptime returns validator's epoch uptime in nanoseconds.
func (s *PublicAbftAPI) GetEpochUptime(ctx context.Context, validatorID hexutil.Uint) (hexutil.Uint64, error) {
	v, err := s.b.GetUptime(ctx, idx.ValidatorID(validatorID))
	if err != nil {
		return 0, err
	}
	if v == nil {
		return 0, nil
	}
	return hexutil.Uint64(v.Uint64()), nil
}

// GetOriginatedEpochFee returns validator's originated epoch fee.
func (s *PublicAbftAPI) GetOriginatedEpochFee(ctx context.Context, validatorID hexutil.Uint) (*hexutil.Big, error) {
	v, err := s.b.GetOriginatedFee(ctx, idx.ValidatorID(validatorID))
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(v), nil
}
