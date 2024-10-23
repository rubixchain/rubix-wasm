# Bidding contract

A simple contract which takes a Bid amount and stores it when the provided Bid amount is larger than the current Bid amount

## Contract Functions

1. **`place_bid`**

Input Parameters: 
    - `bid_amount`: Amount user is billing to bid

Validations:
    - The input `bid_amount` its value must be more than 3.0
    - If the input `bid_amount` is greater than the current Bid amount (maintained in `dapp/state/bid_state.json`), the current bid is updated with the user provided bid.

Output:
    - Returns a String message informing whethere the state was updated or not 