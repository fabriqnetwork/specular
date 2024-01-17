import { BigNumber } from "ethers";
import { Static, Type } from "@sinclair/typebox";

export type CrossDomainMessage = {
  version: BigNumber;
  nonce: BigNumber;
  sender: string;
  target: string;
  value: BigNumber;
  gasLimit: BigNumber;
  data: string;
};

export const FundDepositRequestBody = Type.Object({
  version: Type.String(),
  nonce: Type.String(),
  sender: Type.String(),
  target: Type.String(),
  value: Type.String(),
  gasLimit: Type.String(),
  data: Type.String(),
  depositHash: Type.String(),
});

export type FundDepositRequestBodyType = Static<typeof FundDepositRequestBody>;

export const FundDepositReplyBody = Type.Object({
  txHash: Type.String(),
});

export type FundDepositReplyBodyType = Static<typeof FundDepositReplyBody>;

export const FundDepositReplyErrorBody = Type.Object({
  error: Type.String(),
});

export type FundDepositReplyErrorBodyType = Static<
  typeof FundDepositReplyErrorBody
>;
