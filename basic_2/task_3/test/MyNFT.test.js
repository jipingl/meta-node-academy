const { expect } = require("chai");
const { getNamedAccounts, deployments, ethers } = require("hardhat");

const setup = deployments.createFixture(async () => {
  // 部署合约
  const { MyNFT } = await deployments.fixture(["MyNFT"]);
  // 获取合约实例
  const MyNFTContract = await ethers.getContractAt("MyNFT", MyNFT.address);
  // 获取账户
  const { deployer, owner } = await getNamedAccounts();
  // 返回合约实例和账户
  return {
    MyNFT: MyNFTContract,
    deployer,
    owner,
  };
});

describe("MyNFT", function () {
  it("should have correct contract name", async function () {
    const { MyNFT } = await setup();
    expect(await MyNFT.name()).to.equal("MyNFT");
    expect(await MyNFT.symbol()).to.equal("MNFT");
  });

  it("should mint NTF to owner", async function () {
    const { MyNFT, owner } = await setup();
    await MyNFT.mint(owner);
    expect(await MyNFT.ownerOf(0)).to.equal(owner);
  });

  it("should not allow non-owner to mint", async function () {
    const { MyNFT, owner } = await setup();
    const ownerSigner = await ethers.getSigner(owner);
    await expect(MyNFT.connect(ownerSigner).mint(owner)).to.be.revertedWith(
      "Ownable: caller is not the owner"
    );
  });

  it("should allow owner to transfer NFT", async function () {
    const { MyNFT, deployer, owner } = await setup();
    await MyNFT.mint(owner);
    const ownerSigner = await ethers.getSigner(owner);
    await MyNFT.connect(ownerSigner).transferFrom(owner, deployer, 0);
    expect(await MyNFT.ownerOf(0)).to.equal(deployer);
  });
});
