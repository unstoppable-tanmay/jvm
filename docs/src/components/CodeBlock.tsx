"use client";

import { useState } from "react";
import { Check, Copy } from "lucide-react";

interface CodeBlockProps {
  code: string;
}

export function CodeBlock({ code }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative rounded-xl bg-zinc-50 border border-zinc-200/80 p-3.5 font-mono text-sm text-zinc-800 transition-colors hover:border-zinc-300">
      {/* Top-Right Pinned Copy Button */}
      <button
        onClick={handleCopy}
        aria-label="Copy to clipboard"
        className="absolute top-2.5 right-2.5 p-1.5 rounded-md text-zinc-400 hover:text-zinc-900 hover:bg-zinc-200/50 transition-colors z-10 cursor-pointer"
      >
        {copied ? (
          <Check className="w-3.5 h-3.5 text-zinc-900" />
        ) : (
          <Copy className="w-3.5 h-3.5" />
        )}
      </button>

      {/* Code with hidden scrollbar to eliminate visual overlap */}
      <code className="no-scrollbar overflow-x-auto pr-10 select-all whitespace-pre leading-relaxed block">
        {code}
      </code>
    </div>
  );
}
