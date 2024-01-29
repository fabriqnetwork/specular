import {
    Provider,
} from '@ethersproject/abstract-provider'
import { JsonRpcProvider } from '@ethersproject/providers';
import { ethers } from 'ethers';
import { Wallet, BigNumber, Contract } from "ethers";


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

export function delay(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

export function getStorageKey(messageHash: string) {
    return ethers.utils.keccak256(
        ethers.utils.defaultAbiCoder.encode(
            ["bytes32", "uint256"],
            [messageHash, 0],
        ),
    );
}


//change to waitUntilStateBlock
/**
 * Blocking function that only exits once the L1 state root has been relayed by L1Oracle
 * Note that it is possible for the oracle to skip L1 blocks, in this case this function never exits
 * @param {Contract} l1Oracle - the oracle contract deployed on L2
 * @param {string} stateRoot - the state root we are waiting for
 */
export async function waitUntilStateRoot(
    l1Oracle: Contract,
    stateRoot: string,

) {
    console.log(`Waiting for L2 state root ${stateRoot}...`);

    let oracleStateRoot = await l1Oracle.stateRoot();
    while (oracleStateRoot !== stateRoot) {
        oracleStateRoot = await l1Oracle.stateRoot();
        await delay(500);
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
