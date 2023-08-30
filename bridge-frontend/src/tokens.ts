
interface Token {
  l1TokenName: string;
  l1TokenSymbol: string;
  l1TokenContract: string;
  l2TokenName: string;
  l2TokenSymbol: string;
  l2TokenContract: string;
}

const TOKEN: Record<string, Token> = {
  '2': {
    l1TokenName: "Chiado xDai",
    l1TokenSymbol: "xDai",
    l1TokenContract: "",
    l2TokenName: "Specular ETH",
    l2TokenSymbol: "ETH",
    l2TokenContract: "",
  },
  '1': {
    l1TokenName: "Chiado TT",
    l1TokenSymbol: "TT",
    l1TokenContract: "0xEaa45C3fF72eE58FdB13401586fAB905c507F1BE",
    l2TokenName: "Specular TT",
    l2TokenSymbol: "TT",
    l2TokenContract: "0x6B2031b6519268e623CA05F3683708Ed6C6F89df",
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
