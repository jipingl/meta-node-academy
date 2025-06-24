require("dotenv").config();
require("hardhat-deploy");
require("@nomicfoundation/hardhat-toolbox");

const { CHAIN_ID, PRIVATE_KEY, RPC_URL, ETHERSCAN_API_KEY } = process.env;

module.exports = {
  defaultNetwork: "hardhat",
  networks: {
    sepolia: {
      url: RPC_URL,
      accounts: [PRIVATE_KEY],
    },
    hardhat: {
      chainId: 31337,
    },
  },
  solidity: "0.8.28",
  namedAccounts: {
    deployer: {
      default: 0,
    },
  },
  etherscan: {
    apiKey: ETHERSCAN_API_KEY,
  },
};
