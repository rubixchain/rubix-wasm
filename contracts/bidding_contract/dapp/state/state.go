package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const STATE_FILE_NAME = "bid_state.json"


func getStateFilePath() string {
	absFilePath, _ := filepath.Abs(filepath.Join("state", STATE_FILE_NAME))
	return absFilePath
}


type BidState struct {
	Name       string  `json:"name"`
	Project    string  `json:"project"`
	CurrentBid float64 `json:"current_bid"`
}

func LoadState() (*BidState, error) {
	stateFile, err := os.ReadFile(getStateFilePath())
	if err != nil {
		return nil, fmt.Errorf("failed to load state file: %v, err: %v", getStateFilePath(), err)
	}

	var bidState *BidState
	if err := json.Unmarshal(stateFile, &bidState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal BidState struct, err: %v", err)
	}

	return bidState, nil
}

func SaveState(bidState *BidState) error {
	marshalledBidState, err := json.Marshal(bidState)
	if err != nil {
		return fmt.Errorf("failed to marshal BidState struct, err: %v", err)
	}

	return os.WriteFile(getStateFilePath(), marshalledBidState, 0755)
}

func SaveIncomingBid(bidAmount float64) error {
	bidState, err := LoadState()
	if err != nil {
		return fmt.Errorf("saveIncomingBid: %v", err)
	}

	bidState.CurrentBid = bidAmount

	return SaveState(bidState)
}

func GetCurrentBid() (float64, error) {
	bidState, err := LoadState()
	if err != nil {
		return 0, fmt.Errorf("getCurrentBid: %v", err)
	}

	return bidState.CurrentBid, nil
}
