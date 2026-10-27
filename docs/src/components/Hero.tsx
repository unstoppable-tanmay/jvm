"use client";

import { useEffect, useState } from "react";
import { motion } from "motion/react";
import { ExternalLink } from "lucide-react";
import { CodeBlock } from "@/components/CodeBlock";

export function Hero() {
  const [platform, setPlatform] = useState<"mac" | "brew" | "win">("mac");
  const [latestVersion, setLatestVersion] = useState("v0.1.0");

  useEffect(() => {
    fetch("https://api.github.com/repos/unstoppable-tanmay/jvm/releases/latest")
      .then((res) => (res.ok ? res.json() : null))
      .then((data) => {
        if (data?.tag_name) {
          setLatestVersion(data.tag_name);
        }
      })
      .catch(() => { });
  }, []);

  const commands = {
    mac: "curl -fsSL https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.sh | bash",
    brew: "brew tap unstoppable-tanmay/jvm && brew install jvm",
    win: "irm https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.ps1 | iex",
  };

  return (
    <section className="pt-20 pb-16 max-w-4xl mx-auto px-6 text-center">
      {/* Version badge */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
        className="inline-flex items-center gap-2 text-xs font-mono text-zinc-600 mb-6 bg-zinc-100 hover:bg-zinc-200/70 transition-colors px-3 py-1 rounded-full cursor-pointer"
        onClick={() => window.open("https://github.com/unstoppable-tanmay/jvm/releases", "_blank")}
      >
        <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
        <span>jvm {latestVersion}</span>
        <span className="text-zinc-400">•</span>
        <span className="text-zinc-500">Latest Release</span>
      </motion.div>

      {/* Hero Headline */}
      <motion.h1
        initial={{ opacity: 0, y: 12 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4, delay: 0.1 }}
        className="text-4xl sm:text-6xl font-semibold tracking-tight text-zinc-900 leading-[1.1]"
      >
        Java version management, <br className="hidden sm:inline" />
        made simple.
      </motion.h1>

      {/* Value-driven subtitle focused on developer pain points */}
      <motion.p
        initial={{ opacity: 0, y: 12 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4, delay: 0.15 }}
        className="mt-5 text-base sm:text-lg text-zinc-600 max-w-xl mx-auto leading-relaxed"
      >
        Stop wrestling with broken <code className="text-zinc-800 font-mono text-sm bg-zinc-100 px-1.5 py-0.5 rounded">JAVA_HOME</code> paths,
        slow shell startups, and manual downloads. Switch between Java 8, 17, 21, and beyond in milliseconds.
      </motion.p>

      {/* Install block with simple pill tabs & release button */}
      <motion.div
        initial={{ opacity: 0, y: 14 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4, delay: 0.2 }}
        id="install"
        className="mt-10 max-w-xl mx-auto text-left"
      >
        <div className="flex items-center justify-between gap-2 mb-2.5 flex-wrap ">
          <div className="flex items-center gap-1 [&>button]:cursor-pointer">
            <button
              onClick={() => setPlatform("mac")}
              className={`px-3 py-1 text-xs font-medium rounded-lg transition-colors ${platform === "mac"
                  ? "bg-zinc-900 text-white"
                  : "text-zinc-500 hover:text-zinc-900"
                }`}
            >
              macOS / Linux
            </button>
            <button
              onClick={() => setPlatform("brew")}
              className={`px-3 py-1 text-xs font-medium rounded-lg transition-colors ${platform === "brew"
                  ? "bg-zinc-900 text-white"
                  : "text-zinc-500 hover:text-zinc-900"
                }`}
            >
              Homebrew
            </button>
            <button
              onClick={() => setPlatform("win")}
              className={`px-3 py-1 text-xs font-medium rounded-lg transition-colors ${platform === "win"
                  ? "bg-zinc-900 text-white"
                  : "text-zinc-500 hover:text-zinc-900"
                }`}
            >
              Windows
            </button>
          </div>

          <a
            href="https://github.com/unstoppable-tanmay/jvm/releases"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-mono text-zinc-600 hover:text-zinc-900 hover:bg-zinc-100 rounded-lg border border-zinc-200 transition-colors"
          >
            <span>Releases ({latestVersion})</span>
            <ExternalLink className="w-3 h-3 text-zinc-400" />
          </a>
        </div>

        <CodeBlock code={commands[platform]} />
      </motion.div>
    </section>
  );
}
