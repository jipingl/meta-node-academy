module.exports = async ({ getNamedAccounts, deployments, ethers }) => {
  const { deploy } = deployments;
  const { deployer } = await getNamedAccounts();

  const liquidityReceiver = deployer;
  const marketingReceiver = deployer;
  const burnReceiver = ethers.ZeroAddress;

  // uniswap sepolia router address
  const routerAddress = "0xeE567Fe1712Faf6149d80dA1E6934E354124CfE3";

  await deploy("MyToken", {
    from: deployer,
    args: [routerAddress, liquidityReceiver, marketingReceiver, burnReceiver],
    log: true,
  });
};

module.exports.tags = ["MyToken"];
