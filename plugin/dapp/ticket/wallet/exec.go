// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"github.com/33cn/chain33/types"
	ty "github.com/33cn/plugin/plugin/dapp/ticket/types"
)

func (policy *ticketPolicy) On_CloseTickets(req *types.ReqNil) (types.Message, error) {
	operater := policy.getWalletOperate()
	reply, err := policy.forceCloseTicket(operater.GetBlockHeight() + 1)
	if err != nil {
		bizlog.Error("onCloseTickets", "forceCloseTicket error", err.Error())
	} else {
		go func() {
			if len(reply.Hashes) > 0 {
				operater.WaitTxs(reply.Hashes)
				FlushTicket(policy.getAPI())
			}
		}()
	}
	return reply, err
}

func (policy *ticketPolicy) On_WalletGetTickets(req *types.ReqNil) (types.Message, error) {
	tickets, privs, err := policy.getTicketsByStatus(1)
	tks := &ty.ReplyWalletTickets{tickets, privs}
	return tks, err
}

func (policy *ticketPolicy) On_WalletAutoMiner(req *ty.MinerFlag) (types.Message, error) {
	policy.store.SetAutoMinerFlag(req.Flag)
	policy.setAutoMining(req.Flag)
	FlushTicket(policy.getAPI())
	return &types.Reply{IsOk: true}, nil
}