"use client";

import Link from "next/link";
import Image from "next/image";
import { Coffee } from "lucide-react";

export function Navbar() {
  return (
    <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-zinc-100">
      <div className="max-w-5xl mx-auto px-6 h-14 flex items-center justify-between">
        {/* Left: Text Logo & Internal Navigation */}
        <div className="flex items-center gap-8">
          <Link href="/" className="font-semibold tracking-tight text-zinc-900 text-base">
            jvm
          </Link>

          <nav className="hidden sm:flex items-center gap-6 text-xs sm:text-sm font-medium text-zinc-500">
            <a href="#install" className="hover:text-zinc-900 transition-colors">
              Install
            </a>
            <a href="#highlights" className="hover:text-zinc-900 transition-colors">
              Highlights
            </a>
            <a href="#usage" className="hover:text-zinc-900 transition-colors">
              Commands
            </a>
            <a href="#features" className="hover:text-zinc-900 transition-colors">
              Roadmap
            </a>
          </nav>
        </div>

        {/* Right: Social & Actions (Separated with clean divider) */}
        <div className="flex items-center gap-3">
          {/* Buy me a tea / Sponsor */}
          <a
            href="https://github.com/sponsors/unstoppable-tanmay"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center cursor-pointer gap-1.5 px-3 py-1.5 rounded-lg border border-amber-200/80 bg-amber-50/50 text-amber-900 hover:bg-amber-100/60 hover:border-amber-300 transition-colors text-xs font-medium"
            title="Sponsor unstoppable-tanmay"
          >
            <Coffee className="w-3.5 h-3.5 text-amber-700" />
            <span className="hidden sm:inline">Buy me a Tea</span>
          </a>

          <div className="h-4 w-px bg-zinc-200 mx-1 hidden sm:block" />

          {/* Social links */}
          <a
            href="https://github.com/unstoppable-tanmay/jvm"
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs font-medium text-zinc-600 hover:text-zinc-900 transition-colors"
          >
            GitHub
          </a>
          <a
            href="https://linkedin.com/in/unstoppable-tanmay"
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs font-medium text-zinc-600 hover:text-zinc-900 transition-colors"
          >
            LinkedIn
          </a>

          {/* GitHub Avatar */}
          <a
            href="https://tanmaykumarpanda.in"
            target="_blank"
            rel="noopener noreferrer"
            title="Tanmay's Profile"
            className="ml-1 shrink-0 rounded-full ring-1 ring-zinc-200 hover:ring-zinc-400 transition-all overflow-hidden p-0.5"
          >
            <Image
              src="/me.png"
              alt="Tanmay"
              width={24}
              height={24}
              className="rounded-full object-cover"
              unoptimized
            />
          </a>
        </div>
      </div>
    </header>
  );
}
