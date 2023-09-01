
interface Token {
  l1TokenName: string;
  l1TokenSymbol: string;
  l1TokenContract: string;
  l2TokenName: string;
  l2TokenSymbol: string;
  l2TokenContract: string;
}

const TOKEN: Record<string, Token> = {
  '1': {
    l1TokenName: "Chiado xDai",
    l1TokenSymbol: "xDai",
    l1TokenContract: "",
    l2TokenName: "Specular ETH",
    l2TokenSymbol: "ETH",
    l2TokenContract: "",
  },
  '2': {
    l1TokenName: "Chiado TestToken",
    l1TokenSymbol: "TT",
    l1TokenContract: "0x6d014319E0F36651997697C98Da594c7Cf235fa4",
    l2TokenName: "Specular TestToken",
    l2TokenSymbol: "TT",
    l2TokenContract: "0x6A358FD7B7700887b0cd974202CdF93208F793E2",
  }
};

const erc20Abi = [
  "function balanceOf(address owner) view returns (uint256)",
  "function totalSupply() view returns (uint256)",
  "function approve(address spender, uint256 value) returns (bool)",
  "function allowance(address owner, address spender) view returns (uint256)",
  "function transfer(address to, uint256 value) returns (bool)",
  "function transferFrom(address from, address to, uint256 value) returns (bool)",
  "event Transfer(address indexed from, address indexed to, uint256 value)",
  "event Approval(address indexed owner, address indexed spender, uint256 value)"
]

export {
  TOKEN,
  erc20Abi
};
