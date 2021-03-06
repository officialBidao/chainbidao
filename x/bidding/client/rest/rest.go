package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"bidao/x/bidding"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type placeBidReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	AuctionID string       `json:"auction_id"`
	Bidder    string       `json:"bidder"`
	Bid       string       `json:"bid"`
	Lot       string       `json:"lot"`
}

const (
	restAuctionID = "auction_id"
	restBidder    = "bidder"
	restBid       = "bid"
	restLot       = "lot"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(fmt.Sprintf("/bidding/getauctions"), queryGetAuctionsHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/bidding/bid/{%s}/{%s}/{%s}/{%s}", restAuctionID, restBidder, restBid, restLot), bidHandlerFn(cdc, cliCtx)).Methods("PUT")
}

func queryGetAuctionsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData("/custom/bidding/getauctions", nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func bidHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req placeBidReq
		vars := mux.Vars(r)
		strAuctionID := vars[restAuctionID]
		bechBidder := vars[restBidder]
		strBid := vars[restBid]
		strLot := vars[restLot]

		auctionID, err := bidding.NewIDFromString(strAuctionID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		bidder, err := sdk.AccAddressFromBech32(bechBidder)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		bid, err := sdk.ParseCoin(strBid)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		lot, err := sdk.ParseCoin(strLot)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := bidding.NewMsgPlaceBid(auctionID, bidder, bid, lot)
		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})

	}
}
