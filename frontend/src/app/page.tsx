"use client"
import React from 'react'
import Image from 'next/image'
import starImage from "@/app/star.png"
import { ArrowRight } from 'lucide-react';


const Home: React.FC = () => {
  
    return (
      <>
      <div className="absolute inset-0 -z-10 h-full w-full items-center px-5 py-14 [background:radial-gradient(125%_125%_at_50%_10%,#000_40%,#63e_100%)] overflow-hidden" >
        <div className="relative mx-auto flex max-w-2xl flex-col items-center">
              <div className="mb-8 flex">
              <a
                  href="/dashboard"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex"
                >
                  <span className="relative inline-block overflow-hidden rounded-full p-[1px]">
                    <span className="absolute inset-[-1000%] animate-[spin_2s_linear_infinite] bg-[conic-gradient(from_90deg_at_50%_50%,#a9a9a9_0%,#0c0c0c_50%,#a9a9a9_100%)] dark:bg-[conic-gradient(from_90deg_at_50%_50%,#171717_0%,#737373_50%,#171717_100%)]" />
                    <div className="inline-flex h-full w-full cursor-pointer justify-center rounded-full bg-white px-3 py-1 text-xs font-medium leading-5 text-slate-600 backdrop-blur-xl dark:bg-black dark:text-slate-200">
                      Grab WPTX ⚡️
                      <span className="inline-flex items-center pl-2 text-black dark:text-white">
                        Now{' '}
                        <ArrowRight
                          className="pl-0.5 text-black dark:text-white"
                          size={16}
                        />
                      </span>
                    </div>
                  </span>
                  </a>
              </div>
        </div>
              <h2 className="text-center text-3xl font-medium text-gray-900 dark:text-gray-50 sm:text-6xl md: mt-10">
                Airdrop for{' '}
                <span className="animate-text-gradient inline-flex bg-gradient-to-r from-neutral-900 via-slate-500 to-neutral-500 bg-[100%_auto] bg-clip-text leading-tight text-transparent dark:from-neutral-100 dark:via-slate-400 dark:to-neutral-400">
                  WPTX token
                </span>
              </h2>
              <p className="mt-6 text-center text-lg leading-6 text-gray-600 dark:text-gray-200">
                Wiptrax Token is Brand new Token made on{' '}
                <span className="cursor-wait opacity-70">Kalp DLT Chain</span>
              </p>
              <div className=' relative flex items-center justify-center -z-10 mt-8 md:mt-10'>
              <Image
              src={starImage}
              width={500}
              height={500}
              alt="Picture of the author"
              />
              </div>
      </div>
    </>
    )
  }

export default Home;