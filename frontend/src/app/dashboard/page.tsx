"use client"
import React, { useEffect, useState } from 'react'
import { useKalpApi } from '@/hooks/useKalpAPI'

const Page = () => {

 
  const { claim, balanceOf, totalSupply, transferFrom ,loading } = useKalpApi();
  const [walletAddress, setWalletAddress] = useState("");
  const [balance, setBalance] = useState(0);
  const [totalAirdrop, setTotalAirdrop] = useState(0);
  const [fromtrx, setFromTrx] = useState("");
  const [totrx, setToTrx] = useState("");
  const [valuetrx, setValueTrx] = useState(1);

  const handleClaim = async () => {
    try {
      const data = await claim(walletAddress);
      await handleTotalSupply();
      console.log('Claim successful:', data);
    } catch (err) {
      console.error('Claim error:', err);
    }
  };

  const handleBalanceOf = async () => {
    try {
      const data = await balanceOf(walletAddress);
      setBalance(data.result.result)
      console.log('Balance:', data);
    } catch (err) {
      console.error('BalanceOf error:', err);
    }
  };

  const handleTransferFrom = async () => {
    try {
      const data = await transferFrom(fromtrx, totrx , valuetrx);
      console.log('transfer:', data);
    } catch (err) {
      console.error('BalanceOf error:', err);
    }
  };

  const handleTotalSupply = async () => {
    try {
      const data = await totalSupply();
      setTotalAirdrop(data.result.result)
      console.log('Total Supply:', data);
    } catch (err) {
      console.error('TotalSupply error:', err);
    }
  };


  useEffect(() => {
    handleTotalSupply()
  }, [handleClaim]);

  return (
    <div className="absolute inset-0 -z-10 h-screen w-full items-center px-5 py-10 [background:radial-gradient(125%_125%_at_50%_10%,#000_40%,#63e_100%)]">
    <div className='flex flex-col justify-center items-center'>

      {/* <div className='border-2 py-4 px-16 mt-8 rounded-lg text-4xl w-fit'>Airdrop Machine</div> */}
      <div className='lg:flex gap-12'>
      <div className='flex flex-col border-2 p-8 mt-8 rounded-lg w-fit gap-3'>
        Enter Your Address To Calim :
        <input placeholder='Enter your wallet address' type="text" className='border p-2 rounded-lg w-56 text-black' onChange={(e) => setWalletAddress(e.target.value)} />
        <button className='p-2 rounded-lg bg-[#63e] hover:bg-[#8059eb] text-white disabled:bg-blue-400' onClick={handleClaim} disabled={loading}>{loading ? "Please wait.. " : "Claim"}</button>
      </div>
      <div className='flex flex-col border-[1px] p-8 mt-8 rounded-lg w-fit gap-3'>
          Total WPTX Token Claimed :
          <p className='text-6xl text-[#63e] font-bold w-56'>{totalAirdrop}</p>
        </div>
      </div>

      <div className='lg:flex gap-12'>
        <div className='flex flex-col border-2 p-8 mt-8 rounded-lg w-fit gap-3'>
          My Balance :
          <input placeholder='Enter your wallet address' type="text" className='border p-2 rounded-lg w-56 text-black' onChange={(e) => setWalletAddress(e.target.value)} />
          <button className=' p-2 rounded-lg bg-[#63e] hover:bg-[#8059eb] text-white' onClick={handleBalanceOf}>See</button>

          <p className='text-2xl font-bold w-56'>Balance: <span className='text-[#63e] text-4xl'> {balance}</span></p>
        </div>
        <div className='flex flex-col border-2 p-8 mt-8 rounded-lg w-fit gap-3'>
          Transfer WPTX token :
          <input placeholder='FROM' type="text" className='border p-2 rounded-lg w-56 text-black' onChange={(e) => setFromTrx(e.target.value)} />
          <input placeholder='To' type="text" className='border p-2 rounded-lg w-56 text-black' onChange={(e) => setToTrx(e.target.value)} />
          <input placeholder='No of Token' type="number" className='border p-2 rounded-lg w-56 text-black' onChange={(e) => setValueTrx(Number(e.target.value))} />
          <button className=' p-2 rounded-lg bg-[#63e] hover:bg-[#8059eb] text-white' onClick={handleTransferFrom}>Transfer</button>
        </div>

      </div>

    </div>
    </div>
  )
}

export default Page