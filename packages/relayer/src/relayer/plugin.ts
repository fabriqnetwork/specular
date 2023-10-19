import { FastifyPluginAsync, RouteHandlerMethod, RouteHandler } from "fastify";
import { FastifyRequest, FastifyReply } from "fastify";
import { FastifyPluginAsyncTypebox } from "@fastify/type-provider-typebox";
import cors from '@fastify/cors';

import {
  FundDepositRequestBodyType,
  FundDepositRequestBody,
  RelayerService,
  FundDepositReplyBodyType,
  FundDepositReplyBody,
  FundDepositReplyErrorBody,
  FundDepositReplyErrorBodyType,
  CrossDomainMessage,
} from ".";
import { BigNumber, ethers } from "ethers";

declare module "fastify" {
  interface FastifyInstance {
    relayer: RelayerService;
  }
}

export interface RelayerPluginOptions {
  relayer: RelayerService;
}

export const relayerPlugin: FastifyPluginAsyncTypebox<
  RelayerPluginOptions
> = async (fastify, opts) => {
  fastify.decorate("relayer", opts.relayer);

  fastify.register(cors, {
    origin: '*',
    methods: ['POST'],
  });

  fastify.post<{
    Body: FundDepositRequestBodyType;
    Reply: FundDepositReplyBodyType | FundDepositReplyErrorBodyType;
  }>(
    "/fundDeposit",
    {
      schema: {
        body: FundDepositRequestBody,
        response: {
          200: FundDepositReplyBody,
          400: FundDepositReplyErrorBody,
          500: FundDepositReplyErrorBody,
        },
      },
    },
    async function (request, reply) {
      let {
        nonce: nonceStr,
        sender,
        target,
        value: valueStr,
        gasLimit: gasLimitStr,
        data,
        depositHash,
      } = request.body;
      let nonce: BigNumber;
      try {
        nonce = BigNumber.from(nonceStr);
      } catch (e) {
        reply.status(400).send({ error: "Invalid nonce" });
        return;
      }
      let value: BigNumber;
      try {
        value = BigNumber.from(valueStr);
      } catch (e) {
        reply.status(400).send({ error: "Invalid value" });
        return;
      }
      let gasLimit: BigNumber;
      try {
        gasLimit = BigNumber.from(gasLimitStr);
      } catch (e) {
        reply.status(400).send({ error: "Invalid gas limit" });
        return;
      }
      if (!ethers.utils.isAddress(sender)) {
        reply.status(400).send({ error: "Invalid sender address" });
        return;
      }
      if (!ethers.utils.isAddress(target)) {
        reply.status(400).send({ error: "Invalid target address" });
        return;
      }
      if (data != "0x") {
        reply.status(400).send({ error: "Invalid data" });
        return;
      }
      if (!ethers.utils.isHexString(depositHash, 32)) {
        reply.status(400).send({ error: "Invalid deposit hash" });
        return;
      }
      const targetBalance = await this.relayer.l2Provider.getBalance(target);
      console.log("targetBalance", targetBalance.toString());
      console.log("depositFundingThreshold", this.relayer.config.depositFundingThreshold);
      if (targetBalance.gt(this.relayer.config.depositFundingThreshold)) {
        reply.status(400).send({ error: "Target balance too high" });
        return;
      }
      const depositTx: CrossDomainMessage = {
        nonce,
        sender,
        target,
        value,
        gasLimit,
        data,
      };
      try {
        const txHash = await this.relayer.fundDeposit(
          depositTx,
          depositHash,
          this.relayer.state.lastUpdatedL1OracleBlockNumber,
        );
        request.log.info(`funded deposit ${txHash}`);
        reply.status(200).send({ txHash });
      } catch (e) {
        let message = "Unknown error";
        if (e instanceof Error) {
          message = e.message;
        }
        request.log.error(`failed to fund deposit: ${message}`);
        reply.status(500).send({ error: message });
      }
    }
  );
};
