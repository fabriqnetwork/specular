# RollupBase
[Git Source](https://github.com/SpecularL2/specular/blob/c54213cfb14aca9d44e839341f672dd978834f68/src/Rollup.sol)

**Inherits:**
[IRollup](/src/IRollup.sol/contract.IRollup.md), Initializable, UUPSUpgradeable, OwnableUpgradeable


## State Variables
### confirmationPeriod

```solidity
uint256 public confirmationPeriod;
```


### challengePeriod

```solidity
uint256 public challengePeriod;
```


### minimumAssertionPeriod

```solidity
uint256 public minimumAssertionPeriod;
```


### maxGasPerAssertion

```solidity
uint256 public maxGasPerAssertion;
```


### baseStakeAmount

```solidity
uint256 public baseStakeAmount;
```


### vault

```solidity
address public vault;
```


### stakeToken

```solidity
IERC20 public stakeToken;
```


### sequencerInbox

```solidity
ISequencerInbox public sequencerInbox;
```


### assertions

```solidity
AssertionMap public override assertions;
```


### verifier

```solidity
IVerifier public verifier;
```


## Functions
### __RollupBase_init


```solidity
function __RollupBase_init() internal onlyInitializing;
```

## Structs
### Staker

```solidity
struct Staker {
    bool isStaked;
    uint256 amountStaked;
    uint256 assertionID;
    address currentChallenge;
}
```

### Zombie

```solidity
struct Zombie {
    address stakerAddress;
    uint256 lastAssertionID;
}
```

