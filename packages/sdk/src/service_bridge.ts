import {
  Provider,
  TransactionReceipt,
  TransactionResponse,
  TransactionRequest,
} from "@ethersproject/abstract-provider";

import L1OracleABI from "./abis/L1Oracle.json";
import L1PortalABI from "./abis/L1Portal.json";
import L1StandardBridgeABI from "./abis/L1StandardBridge.json";
import L2PortalABI from "./abis/L2Portal.json";
import RollupABI from "./abis/Rollup.json";
import L2StandarBridgeABI from "./abis/L2StandardBridge.json";
import { BridgeTransaction } from "./interfaces/types";
import ERC20 from "@openzeppelin/contracts/build/contracts/ERC20.json";
import { Signer } from "@ethersproject/abstract-signer";
import { ethers, BigNumber, Overrides, Contract } from "ethers";
import {
  AddressLike,
  MessageStatus,
  NumberLike,
  SignerOrProviderLike,
  L2ChainID,
  L1ChainID,
  MessageType,
} from "./interfaces/types";
import { toNumber, toSignerOrProvider } from "./utils";
import {
  getDepositProof,
  getWithdrawalProof,
  waitUntilBlockConfirmed,
  hexlifyBlockNum,
  getCurrentL1Assertion,
} from "./utils/helpers";

import {
  CONTRACT_ADDRESSES,
  DEFAULT_L2_CONTRACT_ADDRESSES,
  getL1ContractsByNetworkId,
} from "./utils/constants";

import {
  L2Portal,
  L2Portal__factory,
  L1Portal,
  L1Portal__factory,
  L1Oracle,
  L1Oracle__factory,
  L1StandardBridge__factory,
  L1StandardBridge,
  L2StandardBridge,
  L2StandardBridge__factory,
  Rollup,
  Rollup__factory,
} from "./types/contracts";
import { l2PortalAddress } from "./constants";
import { JsonRpcProvider } from "@ethersproject/providers";

export class ServiceBridge {
  /**
   * Provider connected to the L1 chain.
   */
  public l1SignerOrProvider: Signer | Provider;

  /**
   * Provider connected to the L2 chain.
   */
  public l2SignerOrProvider: Signer | Provider;

  /**
   * Chain ID for the L1 network.
   */
  public l1ChainId: number;

  /**
   * Chain ID for the L2 network.
   */
  public l2ChainId: number;

  readonly l1Oracle: L1Oracle;
  readonly l2Portal: L2Portal;
  readonly l1Portal: L1Portal;
  readonly l1Rollup: Rollup;
  readonly l1RollupContract: Contract;
  readonly l1StandardBridge: L1StandardBridge;
  readonly l2StandardBridge: L2StandardBridge;

  public l1RPCprovider: JsonRpcProvider;
  public l2RPCprovider: JsonRpcProvider;

  /**
   * Creates a new Messenger instance.
   *
   * @param opts Options for the provider.
   * @param opts.l1SignerOrProvider Signer or Provider for the L1 chain, or a JSON-RPC url.
   * @param opts.l2SignerOrProvider Signer or Provider for the L2 chain, or a JSON-RPC url.
   * @param opts.l1ChainId Chain ID for the L1 chain.
   * @param opts.l2ChainId Chain ID for the L2 chain.
   */
  constructor(opts: {
    l1SignerOrProvider: SignerOrProviderLike;
    l2SignerOrProvider: SignerOrProviderLike;
    l1ChainId: NumberLike;
    l2ChainId: NumberLike;
  }) {
    this.l1SignerOrProvider = toSignerOrProvider(opts.l1SignerOrProvider);
    this.l2SignerOrProvider = toSignerOrProvider(opts.l2SignerOrProvider);

    this.l1RPCprovider = new ethers.providers.JsonRpcProvider(
      "http://localhost:8545",
    );
    this.l2RPCprovider = new ethers.providers.JsonRpcProvider(
      "http://localhost:4011",
    );

    try {
      this.l1ChainId = toNumber(opts.l1ChainId);
    } catch (err) {
      throw new Error(`L1 chain ID is missing or invalid: ${opts.l1ChainId}`);
    }

    try {
      this.l2ChainId = toNumber(opts.l2ChainId);
    } catch (err) {
      throw new Error(`L2 chain ID is missing or invalid: ${opts.l2ChainId}`);
    }

    const L1StandardBridgeAddress = getL1ContractsByNetworkId(
      this.l1ChainId,
    ).L1StandardBridge.toString();
    const L1PortalAddress = getL1ContractsByNetworkId(
      this.l1ChainId,
    ).L1Portal.toString();
    const l1RollupAddress = getL1ContractsByNetworkId(
      this.l1ChainId,
    ).L1Rollup.toString();

    this.l1Rollup = Rollup__factory.connect(
      l1RollupAddress,
      this.l1SignerOrProvider,
    );
    this.l1StandardBridge = L1StandardBridge__factory.connect(
      L1StandardBridgeAddress,
      this.l1SignerOrProvider,
    );
    this.l2StandardBridge = L2StandardBridge__factory.connect(
      DEFAULT_L2_CONTRACT_ADDRESSES.L2StandardBridge.toString(),
      this.l2SignerOrProvider,
    );
    this.l1Portal = L1Portal__factory.connect(
      L1PortalAddress,
      this.l1SignerOrProvider,
    );
    this.l1Oracle = L1Oracle__factory.connect(
      DEFAULT_L2_CONTRACT_ADDRESSES.L1Oracle.toString(),
      this.l2SignerOrProvider,
    );
    this.l2Portal = L2Portal__factory.connect(
      DEFAULT_L2_CONTRACT_ADDRESSES.L2Portal.toString(),
      this.l2SignerOrProvider,
    );
    this.l1RollupContract = new Contract(
      l1RollupAddress,
      this.l1Rollup.interface,
      this.l1RPCprovider,
    );
  }

  async getL1OracleBlockNumber() {
    return await this.l1Oracle.number();
  }
  /**
   * Provider connected to the L1 chain.
   */
  get l1Provider(): Provider | undefined {
    if (this.l1SignerOrProvider) {
      if (Provider.isProvider(this.l1SignerOrProvider)) {
        return this.l1SignerOrProvider;
      } else {
        return this.l1SignerOrProvider.provider;
      }
    } else {
      return undefined;
    }
  }

  /**
   * Provider connected to the L2 chain.
   */
  get l2Provider(): Provider | undefined {
    if (this.l2SignerOrProvider) {
      if (Provider.isProvider(this.l2SignerOrProvider)) {
        return this.l2SignerOrProvider;
      } else {
        return this.l2SignerOrProvider.provider;
      }
    } else {
      return undefined;
    }
  }

  /**
   * Signer connected to the L1 chain.
   */
  get l1Signer(): Signer {
    if (Provider.isProvider(this.l1SignerOrProvider)) {
      throw new Error(`messenger has no L1 signer`);
    } else {
      return this.l1SignerOrProvider;
    }
  }

  /**
   * Signer connected to the L2 chain.
   */
  get l2Signer(): Signer {
    if (Provider.isProvider(this.l2SignerOrProvider)) {
      throw new Error(`messenger has no L2 signer`);
    } else {
      return this.l2SignerOrProvider;
    }
  }

  /**
   * Return the status of the deposit message.
   *
   * @returns MessageStatus
   */
  public async getDepositStatus(
    bridgeTxLogs: TransactionReceipt,
  ): Promise<MessageStatus> {
    const desiredBlockNumber = bridgeTxLogs.blockNumber;
    const currentOracleBlockNumber = await this.l1Oracle.number();

    const depositEvent = this.l1Portal.interface.parseLog(bridgeTxLogs.logs[1]);
    const depositHash = depositEvent.args.depositHash;

    const depositFinalized = await this.l2Portal.finalizedDeposits(depositHash);
    if (depositFinalized) {
      return MessageStatus.DONE;
    } else if (desiredBlockNumber <= currentOracleBlockNumber.toNumber()) {
      return MessageStatus.READY;
    } else {
      return MessageStatus.PENDING;
    }
  }

  /**
   * Return the status of the withdrawal message.
   *
   * @returns MessageStatus
   */
  public async getWithdrawalStatus(
    bridgeTxLogs: TransactionReceipt,
  ): Promise<MessageStatus> {
    const desiredBlockNumber = bridgeTxLogs.blockNumber;
    const [assertionId, assertionBlockNum] = await getCurrentL1Assertion(
      this.l1RollupContract,
    );

    const withdrawEvent = this.l2Portal.interface.parseLog(
      bridgeTxLogs.logs[1],
    );
    const withdrawalHash = withdrawEvent.args.withdrawalHash;

    const withdrawalFinalized =
      await this.l1Portal.finalizedWithdrawals(withdrawalHash);

    if (withdrawalFinalized) {
      return MessageStatus.DONE;
    } else if (assertionBlockNum >= desiredBlockNumber) {
      return MessageStatus.READY;
    } else {
      return MessageStatus.PENDING;
    }
  }
  /**
   * Returns the last X transactions from a specific address on both L1 and L2
   *
   * @returns MessageStatus
   */
  public async getLastTransactionsFromAddress(
    account: string,
    txsNumber: number,
  ): Promise<BridgeTransaction[]> {
    const txs: BridgeTransaction[] = [];

    const latestL1BlockNumber = await this.l1RPCprovider.getBlockNumber();

    // Create filter for the ETHBridgeInitiated event
    const bridgeInitiatedFilter = {
      address: this.l1StandardBridge.address,
      topics: [
        ethers.utils.id("ETHBridgeInitiated(address,address,uint256,bytes)"), // Event signature with parameter types only
        ethers.utils.hexZeroPad(account, 32),
      ],
      fromBlock: 0, // You can set this to a specific block number if you want to start searching from a specific block
      toBlock: latestL1BlockNumber,
    };

    // Fetch logs matching the filter
    const depositBridgeLogs = await this.l1RPCprovider.getLogs(
      bridgeInitiatedFilter,
    );

    for (const log of depositBridgeLogs) {
      const portalDepositFilter = {
        address: this.l2Portal.address,
        topics: [
          ethers.utils.id(
            "DepositInitiated(uint256,address,address,uint256,uint256,bytes,bytes32)",
          ),
          ethers.utils.hexZeroPad(account, 32),
        ],
        fromBlock: 0,
        toBlock: latestL1BlockNumber,
      };

      const portalLogs = await this.l1RPCprovider.getLogs(portalDepositFilter);

      const depositLog = portalLogs.find(
        (p) => p.args.value === log.args.amount,
      );
      const depositHash = depositLog?.args.depositHash;
      const depositAmount = log.args.amount;

      const message = {
        version: 0n,
        nonce: depositLog?.args.nonce as bigint,
        sender: depositLog?.args.sender as AddressLike,
        target: depositLog?.args.target as AddressLike,
        value: depositLog?.args.value as bigint,
        gasLimit: depositLog?.args.gasLimit as bigint,
        data: depositLog?.args.data,
      };

      txs.push({
        messageHash: depositHash,
        amount: depositAmount,
        block: log.blockNumber,
        type: MessageType.DEPOSIT,
        action: {
          status: MessageStatus.READY,
          message,
          contract: this.l1StandardBridge.address,
          chain: this.l1ChainId,
        },
      });
    }

    const latestL2BlockNumber = await this.l2RPCprovider.getBlockNumber();

    const withdrawalInitiatedFilter = {
      address: this.l2StandardBridge.address,
      topics: [
        ethers.utils.id("ETHBridgeInitiated(address,address,uint256,bytes)"), // Event signature with parameter types only
        ethers.utils.hexZeroPad(account, 32),
      ],
      fromBlock: 0, // You can set this to a specific block number if you want to start searching from a specific block
      toBlock: latestL2BlockNumber,
    };

    // Fetch logs matching the filter
    const withdrawalBridgeLogs = await this.l2RPCprovider.getLogs(
      withdrawalInitiatedFilter,
    );

    for (const log of withdrawalBridgeLogs) {
      const portalWithdrawalFilter = {
        address: this.l1Portal.address,
        topics: [
          ethers.utils.id(
            "WithdrawalInitiated(uint256,address,address,uint256,uint256,bytes,bytes32)",
          ),
          ethers.utils.hexZeroPad(account, 32),
        ],
        fromBlock: 0,
        toBlock: latestL2BlockNumber,
      };

      const portalLogsWithdrawal = await this.l2RPCprovider.getLogs(
        portalWithdrawalFilter,
      );

      const withdrawalLog = portalLogsWithdrawal.find(
        (p) => p.args.value === log.args.amount,
      );
      const withdrawalHash = withdrawalLog?.args.depositHash;
      const withdrawalAmount = log.args.amount;

      const message = {
        version: 0n,
        nonce: withdrawalLog?.args.nonce as bigint,
        sender: withdrawalLog?.args.sender as AddressLike,
        target: withdrawalLog?.args.target as AddressLike,
        value: withdrawalLog?.args.value as bigint,
        gasLimit: withdrawalLog?.args.gasLimit as bigint,
        data: withdrawalLog?.args.data,
      };

      txs.push({
        messageHash: withdrawalHash,
        amount: withdrawalAmount,
        block: log.blockNumber,
        type: MessageType.WITHDRAWAL,
        action: {
          status: MessageStatus.READY,
          message,
          contract: this.l2StandardBridge.address,
          chain: this.l2ChainId,
        },
      });
    }

    return txs.sort((a, b) => {
      if (a.action.status === MessageStatus.READY) {
        return -1;
      }

      if (
        a.action.status === MessageStatus.PENDING &&
        b.action.status === MessageStatus.DONE
      ) {
        return -1;
      }

      return 1;
    });
  }

  /**
   * Deposits ETH into the L2 chain.
   *
   * @param amount Amount of ETH to deposit (in wei).
   * @param opts Additional options.
   * @param opts.signer Optional signer to use to send the transaction.
   * @param opts.recipient Optional address to receive the funds on L2. Defaults to sender.
   * @param opts.l2GasLimit Optional gas limit to use for the transaction on L2.
   * @returns Transaction response for the deposit transaction.
   */
  public async depositETH(
    amount: NumberLike,
    opts?: {
      recipient?: AddressLike;
      signer?: Signer;
      l2GasLimit?: NumberLike;
    },
  ): Promise<TransactionResponse> {
    return (opts?.signer || this.l1Signer).sendTransaction(
      await this.populateTransaction.depositETH(amount, opts),
    );
  }

  /**
   * Approves spending of a specific token.
   *
   * @param token The L1 or L2 token address.
   * @param amount Amount of the token to approve.
   * @param opts Additional options.
   * @param opts.signer Optional signer to use to send the transaction.
   * @returns Transaction response for the approval transaction.
   */
  public async approveERC20(
    token: AddressLike,
    amount: NumberLike,
    chainId: number,
    opts?: {
      signer?: Signer;
    },
  ): Promise<TransactionResponse> {
    return (opts?.signer || this.l1Signer).sendTransaction(
      await this.populateTransaction.approveERC20(token, amount, chainId, opts),
    );
  }

  /**
   * Deposits ERC20 tokens into the L2 chain.
   *
   * @param l1Token Address of the L1 token.
   * @param l2Token Address of the L2 token.
   * @param amount Amount to deposit.
   * @param opts Additional options.
   * @param opts.signer Optional signer to use to send the transaction.
   * @param opts.recipient Optional address to receive the funds on L2. Defaults to sender.
   * @param opts.l2GasLimit Optional gas limit to use for the transaction on L2.
   * @returns Transaction response for the deposit transaction.
   */
  public async depositERC20(
    l1Token: AddressLike,
    l2Token: AddressLike,
    amount: NumberLike,
    opts?: {
      signer?: Signer;
      l2GasLimit?: NumberLike;
    },
  ): Promise<TransactionResponse> {
    return (opts?.signer || this.l1Signer).sendTransaction(
      await this.populateTransaction.depositERC20(
        l1Token,
        l2Token,
        amount,
        opts,
      ),
    );
  }

  /**
   * Withdraws ETH back to the L1 chain.
   *
   * @param amount Amount of ETH to withdraw.
   * @param opts Additional options.
   * @param opts.signer Optional signer to use to send the transaction.
   * @param opts.recipient Optional address to receive the funds on L1. Defaults to sender.
   * @returns Transaction response for the withdraw transaction.
   */
  public async withdrawETH(
    amount: NumberLike,
    opts?: {
      recipient?: AddressLike;
      signer?: Signer;
    },
  ): Promise<TransactionResponse> {
    return (opts?.signer || this.l2Signer).sendTransaction(
      await this.populateTransaction.withdrawETH(amount, opts),
    );
  }

  /**
   * Withdraws ERC20 tokens back to the L1 chain.
   *
   * @param l1Token Address of the L1 token.
   * @param l2Token Address of the L2 token.
   * @param amount Amount to withdraw.
   * @param opts Additional options.
   * @param opts.signer Optional signer to use to send the transaction.
   * @param opts.recipient Optional address to receive the funds on L1. Defaults to sender.
   * @returns Transaction response for the withdraw transaction.
   */
  public async withdrawERC20(
    l1Token: AddressLike,
    l2Token: AddressLike,
    amount: NumberLike,
    opts?: {
      recipient?: AddressLike;
      signer?: Signer;
    },
  ): Promise<TransactionResponse> {
    return (opts?.signer || this.l2Signer).sendTransaction(
      await this.populateTransaction.withdrawERC20(
        l1Token,
        l2Token,
        amount,
        opts,
      ),
    );
  }

  /**
   * Finalizes the deposit on the L2 chain.
   *
   * @param bridgeTxLogs Deposit transaction receipt.
   * @param opts Additional options.
   * @param opts.signer Optional signer to use to send the transaction.
   * @param opts.l2GasLimit Optional gas limit to use for the transaction on L2.
   *
   * @returns Transaction response for the finalizeDeposit transaction.
   */
  public async finalizeDeposit(
    bridgeTxLogs: TransactionReceipt,
    opts?: {
      signer?: Signer;
      l2GasLimit?: NumberLike;
    },
  ): Promise<TransactionResponse> {
    return (opts?.signer || this.l1Signer).sendTransaction(
      await this.populateTransaction.finalizeDeposit(bridgeTxLogs),
    );
  }

  /**
   * Finalizes the withdrawal on the L1 chain.
   *
   * @param bridgeTxLogs Withdrawal transaction receipt
   * @param opts Additional options.
   * @param opts.signer Optional signer to use to send the transaction.
   * @param opts.l2GasLimit Optional gas limit to use for the transaction on L2.
   */
  public async finalizeWithdrawal(
    bridgeTxLogs: TransactionReceipt,
    opts?: {
      signer?: Signer;
      l2GasLimit?: NumberLike;
    },
  ): Promise<TransactionResponse> {
    return (opts?.signer || this.l1Signer).sendTransaction(
      await this.populateTransaction.finalizeWithdrawal(bridgeTxLogs),
    );
  }

  /**
   * Object that holds the functions that generate transactions to be signed by the user.
   * Follows the pattern used by ethers.js.
   */
  populateTransaction = {
    /**
     * Generates a transaction for depositing some ETH into the L2 chain.
     *
     * @param amount Amount of ETH to deposit.
     * @param opts Additional options.
     * @param opts.recipient Optional address to receive the funds on L2. Defaults to sender.
     * @param opts.l2GasLimit Optional gas limit to use for the transaction on L2.
     * @returns Transaction that can be signed and executed to deposit the ETH.
     */
    depositETH: async (
      amount: NumberLike,
      opts?: {
        recipient?: AddressLike;
        l2GasLimit?: NumberLike;
      },
    ): Promise<TransactionRequest> => {
      const bridgeTx =
        await this.l1StandardBridge.populateTransaction.bridgeETH(200_000, [], {
          value: amount,
        });
      return bridgeTx;
    },

    /**
     * Generates a transaction for approval of spending.
     *
     * @param token L1 or L2 token address.
     * @param amount Amount of the token to approve.
     * @returns Transaction response for the approval transaction.
     */
    approveERC20: async (
      token: AddressLike,
      amount: NumberLike,
      chainId: number,
      opts?: {
        signer?: Signer;
      },
    ): Promise<TransactionRequest> => {
      const l1StandardBridgeAddress = getL1ContractsByNetworkId(
        this.l1ChainId,
      ).L1StandardBridge.toString();
      const l2StandardBridgeAddress =
        DEFAULT_L2_CONTRACT_ADDRESSES.L2StandardBridge;

      // check on which chain to do the approval.
      if (Object.values(L1ChainID).includes(this.l1ChainId)) {
        const tokenContract = new Contract(
          token.toString(),
          ERC20.abi,
          this.l1RPCprovider,
        );
        return tokenContract.populateTransaction.approve(
          this.l1StandardBridge.address,
          amount,
        );
      } else {
        const tokenContract = new Contract(
          token.toString(),
          ERC20.abi,
          this.l2RPCprovider,
        );
        return tokenContract.populateTransaction.approve(
          this.l2StandardBridge.address,
          amount,
        );
      }
    },

    /**
     * Generates a transaction for depositing some ERC20 tokens into the L2 chain.
     *
     * @param l1Token Address of the L1 token.
     * @param l2Token Address of the L2 token.
     * @param amount Amount to deposit.
     * @param opts Additional options.
     * @param opts.recipient Optional address to receive the funds on L2. Defaults to sender.
     * @param opts.l2GasLimit Optional gas limit to use for the transaction on L2.
     * @returns Transaction that can be signed and executed to deposit the tokens.
     */
    depositERC20: async (
      l1Token: AddressLike,
      l2Token: AddressLike,
      amount: NumberLike,
      opts?: {
        recipient?: AddressLike;
        l2GasLimit?: NumberLike;
      },
    ): Promise<TransactionRequest> => {
      const bridgeTx =
        await this.l1StandardBridge.populateTransaction.bridgeERC20(
          l1Token.toString(),
          l2Token.toString(),
          amount,
          200_000,
          [],
        );
      return bridgeTx;
    },

    /**
     * Generates a transaction for withdrawing some ETH back to the L1 chain.
     *
     * @param amount Amount of ETH to withdraw.
     * @param opts Additional options.
     * @param opts.recipient Optional address to receive the funds on L1. Defaults to sender.
     * @returns Transaction that can be signed and executed to withdraw the ETH.
     */
    withdrawETH: async (
      amount: NumberLike,
      opts?: {
        recipient?: AddressLike;
        overrides?: Overrides;
      },
    ): Promise<TransactionRequest> => {
      // check if the withdrawer has the sufficient Ether amount.
      const bridgeTx =
        await this.l2StandardBridge.populateTransaction.bridgeETH(200_000, [], {
          value: amount,
        });

      return bridgeTx;
    },

    /**
     * Generates a transaction for withdrawing some ERC20 tokens back to the L1 chain.
     *
     * @param l1Token Address of the L1 token.
     * @param l2Token Address of the L2 token.
     * @param amount Amount to withdraw.
     * @param opts Additional options.
     * @param opts.recipient Optional address to receive the funds on L1. Defaults to sender.
     * @returns Transaction that can be signed and executed to withdraw the tokens.
     */
    withdrawERC20: async (
      l1Token: AddressLike,
      l2Token: AddressLike,
      amount: NumberLike,
      opts?: {
        recipient?: AddressLike;
      },
    ): Promise<TransactionRequest> => {
      const withdrawalTx =
        await this.l2StandardBridge.populateTransaction.bridgeERC20(
          l2Token.toString(),
          l1Token.toString(),
          amount,
          200_000,
          [],
        );
      return withdrawalTx;
    },

    /**
     * Generates a transaction for finalizing the deposit on the L2 chain.
     *
     * @param bridgeTx Deposit transaction receipt.
     * @returns Transaction that can be signed and executed to finalizethe deposit on L2.
     */
    finalizeDeposit: async (
      bridgeTx: TransactionReceipt,
    ): Promise<TransactionRequest> => {
      const tXwithLogs = await this.l1Provider?.getTransactionReceipt(
        bridgeTx.transactionHash,
      );

      const depositEvent = this.l1Portal.interface.parseLog(bridgeTx.logs[1]);

      const despositMessage = {
        version: 0,
        nonce: depositEvent.args.nonce,
        sender: depositEvent.args.sender,
        target: depositEvent.args.target,
        value: depositEvent.args.value,
        gasLimit: depositEvent.args.gasLimit,
        data: depositEvent.args.data,
      };

      // Get initial block number.
      const initialBlockNumber = bridgeTx.blockNumber;

      // get L1 deposit proof.
      const depositProof = await getDepositProof(
        this.l1Portal.address,
        depositEvent.args.depositHash,
        hexlifyBlockNum(initialBlockNumber),
        this.l1RPCprovider,
      );

      const currentBlockNumber = (await this.l1Oracle.number()).toNumber();
      // If the number of the Oracle block if bigger or equal the transaction is settled on L2.
      if (!(initialBlockNumber <= currentBlockNumber)) {
        throw new Error(
          `The deposit transaction can't be finalized, check it's status.`,
        );
      } else {
        const finalizeTx =
          await this.l2Portal.populateTransaction.finalizeDepositTransaction(
            initialBlockNumber,
            despositMessage,
            depositProof.accountProof,
            depositProof.storageProof,
          );
        return finalizeTx;
      }
    },

    /**
     * Generates a transaction for finalizing the withdrawal on the L1 chain.
     *
     * @param bridgeTx Withdrawal transaction receipt.
     * @returns Transaction that can be signed and executed to finalize the withdrawal on L1.
     */
    finalizeWithdrawal: async (
      bridgeTx: TransactionReceipt,
    ): Promise<TransactionRequest> => {
      const withdrawTxBlockNum = bridgeTx.blockNumber;
      const withdrawEvent = this.l2Portal.interface.parseLog(bridgeTx.logs[1]);

      const withdrawMessage = {
        version: 0,
        nonce: withdrawEvent.args.nonce,
        sender: withdrawEvent.args.sender,
        target: withdrawEvent.args.target,
        value: withdrawEvent.args.value,
        gasLimit: withdrawEvent.args.gasLimit,
        data: withdrawEvent.args.data,
      };
      const withdrawalHash = withdrawEvent.args.withdrawalHash;

      const [assertionId, assertionBlockNum] = await getCurrentL1Assertion(
        this.l1Rollup,
      );

      // Get withdraw proof for the block the assertion committed to.
      const withdrawProof = await getWithdrawalProof(
        this.l2Portal.address,
        withdrawalHash,
        hexlifyBlockNum(assertionBlockNum),
        this.l2RPCprovider,
      );

      // Get block for the block the assertion committed to.
      const rawBlock = await this.l2RPCprovider.send("eth_getBlockByNumber", [
        ethers.utils.hexValue(assertionBlockNum),
        false, // We only want the block header
      ]);
      const l2BlockHash = this.l2RPCprovider.formatter.hash(rawBlock.hash);
      const l2StateRoot = this.l2RPCprovider.formatter.hash(rawBlock.stateRoot);

      const finalizeTx =
        await this.l1Portal.populateTransaction.finalizeWithdrawalTransaction(
          withdrawMessage,
          assertionId,
          l2BlockHash,
          l2StateRoot,
          withdrawProof.accountProof,
          withdrawProof.storageProof,
        );
      return finalizeTx;
    },
  };
}
