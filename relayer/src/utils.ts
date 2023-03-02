import { utils } from "ethers";
import { EncodedBlockHeader } from "./messager";

export function getStorageKey(messageHash: string) {
  return utils.keccak256(
    utils.defaultAbiCoder.encode(["bytes32", "uint256"], [messageHash, 0])
  );
}

export function rawBlockHeaderToEncoded(rawBlock: any): EncodedBlockHeader {
  const parentHash = rawBlock.parentHash;
  const sha3Uncles = rawBlock.sha3Uncles;
  const miner = rawBlock.miner;
  const stateRoot = rawBlock.stateRoot;
  const transactionsRoot = rawBlock.transactionsRoot;
  const receiptsRoot = rawBlock.receiptsRoot;
  const logsBloom = rawBlock.logsBloom;
  const difficulty = utils.hexlify(rawBlock.difficulty, { hexPad: "left" });
  const number = utils.hexlify(utils.hexValue(rawBlock.number), {
    hexPad: "left",
  });
  const gasLimit = utils.hexlify(utils.hexValue(rawBlock.gasLimit), {
    hexPad: "left",
  });
  const gasUsed = utils.hexlify(utils.hexValue(rawBlock.gasUsed), {
    hexPad: "left",
  });
  const timestamp = utils.hexlify(utils.hexValue(rawBlock.timestamp), {
    hexPad: "left",
  });
  const extraData = utils.hexlify(rawBlock.extraData, { hexPad: "left" });
  const mixHash = rawBlock.mixHash;
  const nonce = utils.hexlify(rawBlock.nonce, { hexPad: "left" });
  let header = [
    parentHash,
    sha3Uncles,
    miner,
    stateRoot,
    transactionsRoot,
    receiptsRoot,
    logsBloom,
    difficulty,
    number,
    gasLimit,
    gasUsed,
    timestamp,
    extraData,
    mixHash,
    nonce,
  ];
  header = header.map((v) => {
    if (v == "0x00") return "0x";
    return v;
  });
  return utils.RLP.encode(header);
}
