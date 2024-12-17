package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const STATE_FILE_NAME = "bid_state.json"

type PlaceBidReqState struct {
	BidderDID          string `json:"bidder_did"`
	EncryptedBidAmount string `json:"encrypted_bid_amount"`
}

type BiddingState struct {
	BiddingInfo map[string]string `json:"bidding_info"`
}

func getStateFilePath() string {
	absFilePath, _ := filepath.Abs(filepath.Join("state", STATE_FILE_NAME))
	return absFilePath
}

func GetBidInfo() (*BiddingState, error) {
	return LoadBiddingState()
}

func saveBiddingState(placeBidReqState *BiddingState) error {
	marshalledPlaceBidReqState, err := json.Marshal(placeBidReqState)
	if err != nil {
		return fmt.Errorf("failed to marshal BidState struct, err: %v", err)
	}

	return os.WriteFile(getStateFilePath(), marshalledPlaceBidReqState, 0755)

}

func LoadBiddingState() (*BiddingState, error) {
	// Read the file containing the PlaceBidReqState
	stateFile, err := os.ReadFile(getStateFilePath())
	if err != nil {
		return nil, fmt.Errorf("failed to load bid info file: %v, err: %v", getStateFilePath(), err)
	}

	// Unmarshal JSON data into the PlaceBidReqState struct
	var biddingState *BiddingState
	if err := json.Unmarshal(stateFile, &biddingState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal biddingState struct, err: %v", err)
	}

	return biddingState, nil
}

func SavePlaceBidReqState(req PlaceBidReqState) error {
	biddingState, err := LoadBiddingState()
	if err != nil {
		return fmt.Errorf("error while loading bidding state, err: %v", err)
	}

	if biddingState == nil {
		return fmt.Errorf("bidding state is nil")
	}

	// Save the bidder DID and encrypted bid amount in the state
	biddingState.BiddingInfo[req.BidderDID] = req.EncryptedBidAmount

	return saveBiddingState(biddingState)
}
