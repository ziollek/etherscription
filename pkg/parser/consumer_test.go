package parser

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ziollek/etherscription/pkg/config"
	"github.com/ziollek/etherscription/pkg/model"
	"github.com/ziollek/etherscription/pkg/storage/mocks"
)

func TestShouldConsumeInformationAboutParsedBlockId(t *testing.T) {
	type args struct {
		lastBlock int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Should consume information about parsed block id", args{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			kv := mock_storage.NewMockKVSaver[int](ctrl)
			kv.EXPECT().Set("last_block", tt.args.lastBlock)
			NewStateConsumerService(kv).Consume(tt.args.lastBlock)
		})
	}
}

func TestShouldConsumeParsedTransactions(t *testing.T) {
	type fields struct {
		ttl        time.Duration
		storeAllTx bool
	}
	type args struct {
		transaction       model.Transaction
		subscriptionState map[string]bool
	}
	type expected struct {
		shouldAppendFor []string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			"Should not store transaction that is not subscribed",
			fields{time.Second, false},
			args{model.Transaction{From: "0x1", To: "0x2", Value: 1}, map[string]bool{"0x1": false, "0x2": false}},
			expected{[]string{}},
		},
		{
			"Should store transaction that is subscribed by receiver once",
			fields{time.Second, false},
			args{model.Transaction{From: "0x1", To: "0x2", Value: 1}, map[string]bool{"0x1": false, "0x2": true}},
			expected{[]string{"0x2"}},
		},
		{
			"Should store transaction that is subscribed by sender once",
			fields{time.Second, false},
			args{model.Transaction{From: "0x1", To: "0x2", Value: 1}, map[string]bool{"0x1": true, "0x2": false}},
			expected{[]string{"0x1"}},
		},
		{
			"Should store transaction that is subscribed by sender & receiver twice",
			fields{time.Second, false},
			args{model.Transaction{From: "0x1", To: "0x2", Value: 1}, map[string]bool{"0x1": true, "0x2": true}},
			expected{[]string{"0x1", "0x2"}},
		},
		{
			"Should store transaction that is not subscribed but store all transaction is enabled",
			fields{time.Second, true},
			args{model.Transaction{From: "0x1", To: "0x2", Value: 1}, map[string]bool{"0x1": false, "0x2": false}},
			expected{[]string{"0x1", "0x2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			kv := mock_storage.NewMockKVSaver[string](ctrl)
			for k, exists := range tt.args.subscriptionState {
				kv.EXPECT().Get(k).Return("", exists)
			}
			txStorage := mock_storage.NewMockListSaver[model.Transaction](ctrl)
			for _, address := range tt.expected.shouldAppendFor {
				txStorage.EXPECT().Append(address, tt.args.transaction, tt.fields.ttl)
			}
			s := NewConsumerService(txStorage, kv, &config.StorageConfig{Retention: tt.fields.ttl, StoreAllTransactions: tt.fields.storeAllTx})
			s.Consume(tt.args.transaction)
		})
	}
}
