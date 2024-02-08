
import { JsonRpcProvider } from '@ethersproject/providers';
import { ethers } from 'ethers';
import { BigNumber, Contract } from "ethers";
import {
    Provider,
    TransactionReceipt,
    TransactionResponse,
} from '@ethersproject/abstract-provider'
import { Signer } from '@ethersproject/abstract-signer'
import {
    SignerOrProviderLike,
    ProviderLike,
    TransactionLike,
    NumberLike,
    AddressLike,
} from '../interfaces'


/** 
* Get a depositProof for finalization.
*/
export async function getDepositProof(
    portalAddress: string,
    depositHash: string,
    blockNumber: string,
    l1Provider: JsonRpcProvider
) {
    const proof = await l1Provider.send("eth_getProof", [
        portalAddress,
        [getStorageKey(depositHash)],
        blockNumber,
    ]);

    return {
        accountProof: proof.accountProof,
        storageProof: proof.storageProof[0].proof,
    };
}

/** 
* Get a withdrawalProof for finalization.
*/
export async function getWithdrawalProof(
    portalAddress: string,
    withdrawalHash: string,
    blockNumber: string,
    l2Provider: JsonRpcProvider
) {
    const proof = await l2Provider.send("eth_getProof", [
        portalAddress,
        [getStorageKey(withdrawalHash)],
        blockNumber,
    ]);

    return {
        accountProof: proof.accountProof,
        storageProof: proof.storageProof[0].proof,
    };
}

/** 
* Enforces a milisecond delay. 
*/
export function delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

/** 
* Gets the storage key. 
*/
export function getStorageKey(messageHash: string) {
    return ethers.utils.keccak256(
        ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "uint256"],
            [messageHash, 0],
        ),
    );
}

/**
 * Blocking function that only exits once the block relayed by L1Oracle is >= the blockNumber
 * @param {Contract} l1Oracle - the oracle contract deployed on L2
 * @param {number} blockNumber - the block we are waiting for
 */
export async function waitUntilOracleBlock(
    l1Oracle: Contract,
    blockNumber: number,
) {
    console.log(`Waiting for L1Oracle to relay block #${blockNumber}...`);
    let oracleBlockNumber = await l1Oracle.number();
    while (oracleBlockNumber < blockNumber) {
        await delay(500);
        oracleBlockNumber = await l1Oracle.number();
    }
}

/**
 * Blocking function that only exits once L2 block with blockNum has been confirmed on L1
 * @param {Contract} rollup - the rollup contract deployed on L1
 * @param {number} blockNum - number of the L2 block we are waiting on
 */
export async function waitUntilBlockConfirmed(
    rollup: Contract,
    blockNum: number,
): Promise<[number, number]> {
    let confirmedAssertionId: number | undefined = undefined;
    let confirmedBlockNum: number;

    rollup.on(rollup.filters.AssertionCreated(), () => {
        console.log("AssertionCreated");
    });

    rollup.on(rollup.filters.AssertionConfirmed(), async (id: BigNumber) => {
        const assertionId = id.toNumber();
        const assertion = await rollup.getAssertion(assertionId);
        const assertionBlockNum = assertion.blockNum.toNumber();

        console.log({
            msg: "AssertionConfirmed",
            id: assertionId,
            blockNum: assertionBlockNum,
        });
        if (!confirmedAssertionId && blockNum <= assertionBlockNum) {
            // Found the first assertion to confirm block
            confirmedAssertionId = assertionId;
            confirmedBlockNum = assertionBlockNum;
        }
    });

    console.log(`Waiting for L2 block ${blockNum} to be confirmed...`);
    while (!confirmedAssertionId) {
        await delay(500);
    }

    return [confirmedAssertionId!, confirmedBlockNum!];
}

// eth_getProof block number param cannot have leading zeros
// This function hexlifys blockNum and strips leading zeros
export function hexlifyBlockNum(blockNum: number): string {
    let hexBlockNum = ethers.utils.hexlify(blockNum);
    // Check if the string starts with "0x" and contains more than just "0x".
    if (hexBlockNum.startsWith("0x") && hexBlockNum.length > 2) {
        let strippedString = "0x";

        // Iterate through the characters of the input string starting from the third character (index 2).
        for (let i = 2; i < hexBlockNum.length; i++) {
            if (hexBlockNum[i] !== "0") {
                strippedString += hexBlockNum.substring(i); // Append the remaining characters.
                return strippedString;
            }
        }

        // If all characters are '0', return "0x0".
        return "0x0";
    }

    // If the input is not in the expected format, return it as is.
    return hexBlockNum;
}

export const assert = (condition: boolean, message: string): void => {
    if (!condition) {
        throw new Error(message)
    }
}

/**
 * Converts a SignerOrProviderLike into a Signer or a Provider. Assumes that if the input is a
 * string then it is a JSON-RPC url.
 *
 * @param signerOrProvider SignerOrProviderLike to turn into a Signer or Provider.
 * @returns Input as a Signer or Provider.
 */
export const toSignerOrProvider = (
    signerOrProvider: SignerOrProviderLike
): Signer | Provider => {
    if (typeof signerOrProvider === 'string') {
        return new ethers.providers.JsonRpcProvider(signerOrProvider)
    } else if (Provider.isProvider(signerOrProvider)) {
        return signerOrProvider
    } else if (Signer.isSigner(signerOrProvider)) {
        return signerOrProvider
    } else {
        throw new Error('Invalid provider')
    }
}

/**
 * Converts a ProviderLike into a Provider. Assumes that if the input is a string then it is a
 * JSON-RPC url.
 *
 * @param provider ProviderLike to turn into a Provider.
 * @returns Input as a Provider.
 */
export const toProvider = (provider: ProviderLike): Provider => {
    if (typeof provider === 'string') {
        return new ethers.providers.JsonRpcProvider(provider)
    } else if (Provider.isProvider(provider)) {
        return provider
    } else {
        throw new Error('Invalid provider')
    }
}

/**
 * Pulls a transaction hash out of a TransactionLike object.
 *
 * @param transaction TransactionLike to convert into a transaction hash.
 * @returns Transaction hash corresponding to the TransactionLike input.
 */
export const toTransactionHash = (transaction: TransactionLike): string => {
    if (typeof transaction === 'string') {
        assert(
            ethers.utils.isHexString(transaction, 32),
            'Invalid transaction hash'
        )
        return transaction
    } else if ((transaction as TransactionReceipt).transactionHash) {
        return (transaction as TransactionReceipt).transactionHash
    } else if ((transaction as TransactionResponse).hash) {
        return (transaction as TransactionResponse).hash
    } else {
        throw new Error('Invalid transaction')
    }
}

/**
 * Converts a number-like into an ethers BigNumber.
 *
 * @param num Number-like to convert into a BigNumber.
 * @returns Number-like as a BigNumber.
 */
export const toBigNumber = (num: NumberLike): BigNumber => {
    return ethers.BigNumber.from(num)
}

/**
 * Converts a number-like into a number.
 *
 * @param num Number-like to convert into a number.
 * @returns Number-like as a number.
 */
export const toNumber = (num: NumberLike): number => {
    return toBigNumber(num).toNumber()
}

/**
 * Converts an address-like into a 0x-prefixed address string.
 *
 * @param addr Address-like to convert into an address.
 * @returns Address-like as an address.
 */
export const toAddress = (addr: AddressLike): string => {
    if (typeof addr === 'string') {
        assert(ethers.utils.isAddress(addr), 'Invalid address')
        return ethers.utils.getAddress(addr)
    } else {
        assert(ethers.utils.isAddress(addr.address), 'Invalid address')
        return ethers.utils.getAddress(addr.address)
    }
}
