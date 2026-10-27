"use client";

import { useState } from "react";
import { motion } from "motion/react";
import { CodeBlock } from "@/components/CodeBlock";

export function UsageShowcase() {
  const [activeTab, setActiveTab] = useState<
    "use" | "install" | "ls" | "link" | "unlink" | "setdir" | "current" | "uninstall"
  >("use");

  const commands = {
    use: {
      title: "jvm use <version>",
      desc: "Switches the active JDK version in your current shell tab. Prompts a selector if ambiguous.",
      code: `# Switch immediately in current terminal
jvm use 21

# Switch to specific vendor distribution
jvm use corretto@17

# If multiple versions match, an intuitive selector appears:
# Multiple matching JDKs found for switching active JDK:
#   [1] temurin-21+35   (installed)
#   [2] homebrew-21     (linked -> /opt/homebrew)
# Select an option [1-2]: 1`,
    },
    install: {
      title: "jvm install <version> (alias: i)",
      desc: "Downloads and verifies OpenJDK releases from upstream mirrors or custom URLs.",
      code: `# Install Eclipse Temurin 21 (default vendor)
jvm install 21
# Shorthand alias
jvm i 21

# Install Amazon Corretto 17
jvm i corretto@17

# Install from custom direct URL with SHA256 integrity check
jvm install --from "https://github.com/.../jdk.tar.gz" --sha256 "<hash>" custom-21`,
    },
    ls: {
      title: "jvm ls & jvm ls-remote",
      desc: "Lists all locally installed/linked versions, or queries available remote releases.",
      code: `# List installed and linked JDKs
jvm ls

# List remote releases from Disco API
jvm ls-remote --vendor temurin --version 21
jvm ls-remote --vendor corretto --version 17`,
    },
    link: {
      title: "jvm link <path> [name]",
      desc: "Links an existing external JDK into jvm without downloading duplicate files.",
      code: `# Link macOS native installer bundle (auto-detects release metadata)
jvm link /Library/Java/JavaVirtualMachines/temurin-17.jdk

# Link Homebrew JDK with custom alias
jvm link /opt/homebrew/opt/openjdk@11 homebrew-11`,
    },
    unlink: {
      title: "jvm unlink <name>",
      desc: "Safely unlinks an externally linked JDK without touching original files on disk.",
      code: `# Remove symlink from jvm (external JDK remains untouched)
jvm unlink homebrew-11`,
    },
    setdir: {
      title: "jvm set-dir <path>",
      desc: "Configures a custom JDK storage folder (e.g. external SSD) with migration & linking options.",
      code: `# Change JDK storage directory
jvm set-dir /Volumes/ExternalSSD/jdks

# Prompts:
# Existing JDKs found in ~/.jvm/versions:
#   [1] Move (migrate all existing JDKs to new directory)
#   [2] Link (symlink existing JDKs into new directory)
#   [3] Leave (keep them where they are without linking)`,
    },
    current: {
      title: "jvm current",
      desc: "Displays the active default jvm Java version and system runtime path.",
      code: `jvm current
      
# Output:
# Configured jvm default: temurin-21+35
# Path: /Users/tanmay/.jvm/versions/temurin-21+35`,
    },
    uninstall: {
      title: "jvm uninstall <version> (alias: ui)",
      desc: "Removes a downloaded JDK version from jvm storage.",
      code: `# Uninstall an installed JDK
jvm uninstall 17

# Shorthand alias
jvm ui 17`,
    },
  };

  const tabs = [
    { id: "use", label: "use" },
    { id: "install", label: "install / i" },
    { id: "ls", label: "ls / ls-remote" },
    { id: "link", label: "link" },
    { id: "unlink", label: "unlink" },
    { id: "setdir", label: "set-dir" },
    { id: "current", label: "current" },
    { id: "uninstall", label: "uninstall / ui" },
  ] as const;

  return (
    <section id="usage" className="py-20 max-w-4xl mx-auto px-6">
      <div className="mb-8">
        <h2 className="text-2xl font-semibold tracking-tight text-zinc-900">
          CLI Command Reference
        </h2>
        <p className="text-sm text-zinc-500 mt-1">
          Complete list of commands and shorthand aliases available in jvm.
        </p>
      </div>

      {/* Tabs */}
      <div className="flex items-center gap-1.5 border-b border-zinc-200 pb-2 mb-6 overflow-x-auto code-scroll">
        {tabs.map((t) => (
          <button
            key={t.id}
            onClick={() => setActiveTab(t.id)}
            className={`px-3 py-1.5 text-xs font-mono cursor-pointer font-medium rounded-lg whitespace-nowrap transition-colors ${activeTab === t.id
              ? "bg-zinc-900 text-white"
              : "text-zinc-600 hover:text-zinc-900 hover:bg-zinc-100"
              }`}
          >
            {t.label}
          </button>
        ))}
      </div>

      <motion.div
        key={activeTab}
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.2 }}
        className="space-y-3"
      >
        <div>
          <div className="text-sm font-semibold font-mono text-zinc-900">
            {commands[activeTab].title}
          </div>
          <div className="text-xs text-zinc-500 mt-0.5">
            {commands[activeTab].desc}
          </div>
        </div>
        <CodeBlock code={commands[activeTab].code} />
      </motion.div>
    </section>
  );
}
