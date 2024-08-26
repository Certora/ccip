package withdrawprover

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l2_to_l1_message_passer"
)

func Test_hashLowLevelMessage(t *testing.T) {
	type args struct {
		llm *optimism_l2_to_l1_message_passer.OptimismL2ToL1MessagePasserMessagePassed
	}
	tests := []struct {
		name    string
		args    args
		want    [32]byte
		wantErr bool
	}{
		// tx: https://sepolia-optimism.etherscan.io/tx/0x14e41dac648d2e1c166ca5c44af09c7c1da684b994ae74dc11303b1ac4bd057c
		{
			"success",
			args{
				&optimism_l2_to_l1_message_passer.OptimismL2ToL1MessagePasserMessagePassed{
					Nonce:          decimal.RequireFromString("1766847064778384329583297500742918515827483896875618958121606201292621236").BigInt(),
					Sender:         common.HexToAddress("0x4200000000000000000000000000000000000007"),
					Target:         common.HexToAddress("0x58Cc85b8D04EA49cC6DBd3CbFFd00B4B8D6cb3ef"),
					Value:          big.NewInt(0),
					GasLimit:       big.NewInt(289_160),
					Data:           hexutil.MustDecode("0xD764AD0B00010000000000000000000000000000000000000000000000000000000005700000000000000000000000004200000000000000000000000000000000000010000000000000000000000000FBB0621E0B23B5478B630BD55A5F21F67730B0F10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000C000000000000000000000000000000000000000000000000000000000000001040166A07A0000000000000000000000005589BB8228C07C4E15558875FAF2B859F678D129000000000000000000000000D08A2917653D4E460893203471F0000826FB403400000000000000000000000084F429D9F18D8512FFB2DBD6855DC110C7186F1300000000000000000000000077FFC73ED3B2614D21B3398FE368E989F318B412000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000C00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000"),
					WithdrawalHash: common.HexToHash("0x682D40CCAD4D19BC7A37CB0A6AAEA0E58C2BDEB7F59062AF974A82978ABB9445"),
				},
			},
			common.HexToHash("0x682D40CCAD4D19BC7A37CB0A6AAEA0E58C2BDEB7F59062AF974A82978ABB9445"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashLowLevelMessage(tt.args.llm)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParseMessagePassedLog(t *testing.T) {
	type args struct {
		log *gethtypes.Log
	}
	tests := []struct {
		name    string
		args    args
		want    *optimism_l2_to_l1_message_passer.OptimismL2ToL1MessagePasserMessagePassed
		wantErr bool
	}{
		// tx: https://sepolia-optimism.etherscan.io/tx/0x14e41dac648d2e1c166ca5c44af09c7c1da684b994ae74dc11303b1ac4bd057c
		{
			"success",
			args{
				&gethtypes.Log{
					Topics: []common.Hash{
						common.HexToHash("0x02a52367d10742d8032712c1bb8e0144ff1ec5ffda1ed7d70bb05a2744955054"),
						common.BigToHash(decimal.RequireFromString("1766847064778384329583297500742918515827483896875618958121606201292621236").BigInt()),
						common.HexToHash("0x4200000000000000000000000000000000000007"),
						common.HexToHash("0x58Cc85b8D04EA49cC6DBd3CbFFd00B4B8D6cb3ef"),
					},
					Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000469880000000000000000000000000000000000000000000000000000000000000080682d40ccad4d19bc7a37cb0a6aaea0e58c2bdeb7f59062af974a82978abb94450000000000000000000000000000000000000000000000000000000000000204d764ad0b00010000000000000000000000000000000000000000000000000000000005700000000000000000000000004200000000000000000000000000000000000010000000000000000000000000fbb0621e0b23b5478b630bd55a5f21f67730b0f10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001040166a07a0000000000000000000000005589bb8228c07c4e15558875faf2b859f678d129000000000000000000000000d08a2917653d4e460893203471f0000826fb403400000000000000000000000084f429d9f18d8512ffb2dbd6855dc110c7186f1300000000000000000000000077ffc73ed3b2614d21b3398fe368e989f318b412000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
				},
			},
			&optimism_l2_to_l1_message_passer.OptimismL2ToL1MessagePasserMessagePassed{
				Nonce:          decimal.RequireFromString("1766847064778384329583297500742918515827483896875618958121606201292621236").BigInt(),
				Sender:         common.HexToAddress("0x4200000000000000000000000000000000000007"),
				Target:         common.HexToAddress("0x58Cc85b8D04EA49cC6DBd3CbFFd00B4B8D6cb3ef"),
				Value:          big.NewInt(0),
				GasLimit:       big.NewInt(289_160),
				Data:           hexutil.MustDecode("0xD764AD0B00010000000000000000000000000000000000000000000000000000000005700000000000000000000000004200000000000000000000000000000000000010000000000000000000000000FBB0621E0B23B5478B630BD55A5F21F67730B0F10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000C000000000000000000000000000000000000000000000000000000000000001040166A07A0000000000000000000000005589BB8228C07C4E15558875FAF2B859F678D129000000000000000000000000D08A2917653D4E460893203471F0000826FB403400000000000000000000000084F429D9F18D8512FFB2DBD6855DC110C7186F1300000000000000000000000077FFC73ED3B2614D21B3398FE368E989F318B412000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000C00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000"),
				WithdrawalHash: common.HexToHash("0x682D40CCAD4D19BC7A37CB0A6AAEA0E58C2BDEB7F59062AF974A82978ABB9445"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMessagePassedLog(tt.args.log)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.Data, got.Data)
				require.Equal(t, tt.want.GasLimit, got.GasLimit)
				require.Equal(t, tt.want.Nonce, got.Nonce)
				require.Equal(t, tt.want.Sender, got.Sender)
				require.Equal(t, tt.want.Target, got.Target)
				require.Equal(t, tt.want.Value.String(), got.Value.String())
				require.Equal(t, tt.want.WithdrawalHash, got.WithdrawalHash)
			}
		})
	}
}

func Test_hashMessageHash(t *testing.T) {
	type args struct {
		h [32]byte
	}
	tests := []struct {
		name    string
		args    args
		want    [32]byte
		wantErr bool
	}{
		// tx: https://sepolia-optimism.etherscan.io/tx/0x14e41dac648d2e1c166ca5c44af09c7c1da684b994ae74dc11303b1ac4bd057c
		{
			"success",
			args{
				common.HexToHash("0x682D40CCAD4D19BC7A37CB0A6AAEA0E58C2BDEB7F59062AF974A82978ABB9445"),
			},
			common.HexToHash("0x4ee0474e73e1e4e24aa4502f8f2be6821e35dbe367736ca942566749559b2e34"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashMessageHash(tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashMessageHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hashMessageHash() = %v, want %v", got, tt.want)
			}
		})
	}
}