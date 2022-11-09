#/bin/sh

npx hardhat compile
abigen --abi ./abi/src/AssertionMap.sol/AssertionMap.json --pkg bindings --type AssertionMap --out bindings/AssertionMap.go
abigen --abi ./abi/src/challenge/IChallenge.sol/IChallenge.json --pkg bindings --type IChallenge --out bindings/IChallenge.go
abigen --abi ./abi/src/IRollup.sol/IRollup.json --pkg bindings --type IRollup --out bindings/IRollup.go
abigen --abi ./abi/src/ISequencerInbox.sol/ISequencerInbox.json --pkg bindings --type ISequencerInbox --out bindings/ISequencerInbox.go
