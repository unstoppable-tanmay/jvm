"use client";

import { motion } from "motion/react";
import { SquareCheck, Square } from "lucide-react";

export function FeaturesChecklist() {
  const completed = [
    "Instant version switching in your current shell session without restarting terminal",
    "Auto-switch Java runtime when navigating into projects with .java-version or .jvmrc",
    "Link existing JDKs from Homebrew, Android Studio, and macOS installers without redownloading",
    "Clean unlinking to remove links without touching original external files on disk",
    "Interactive terminal selector when multiple installed or linked versions match a target",
    "Install directly from upstream vendors (Temurin, Corretto, Zulu, GraalVM) via Disco API",
    "Install from custom direct URLs with optional SHA256 integrity verification (--from)",
    "Configurable storage location (e.g. external SSDs) with safe migration prompts (jvm set-dir)",
    "Cross-platform support across macOS (Apple Silicon & Intel), Linux, and Windows (PowerShell & cmd.exe)",
    "Clean interactive uninstallation with prompt to preserve or delete downloaded JDKs",
  ];

  const planned = [
    "Automatic Homebrew formula updates on new GitHub release tags",
    "Built-in proxy and enterprise mirror registry support in user settings",
    "Direct migration wizard to import existing SDKMAN candidate directories",
  ];

  return (
    <section id="features" className="py-20 border-t border-zinc-100">
      <div className="max-w-4xl mx-auto px-6">
        <div className="mb-8">
          <h2 className="text-2xl font-semibold tracking-tight text-zinc-900">
            Features &amp; Capabilities
          </h2>
          <p className="text-sm text-zinc-500 mt-1">
            User-facing capabilities designed to make everyday Java workflows seamless.
          </p>
        </div>

        <div className="space-y-6 pl-1">
          {/* Completed Features */}
          <div className="space-y-3">
            {completed.map((item, i) => (
              <motion.div
                key={item}
                initial={{ opacity: 0, y: 6 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.25, delay: i * 0.02 }}
                className="flex items-start gap-3 text-sm text-zinc-700 leading-snug"
              >
                <SquareCheck className="w-4 h-4 text-zinc-900 shrink-0 mt-0.5" />
                <span>{item}</span>
              </motion.div>
            ))}
          </div>

          {/* Planned Features */}
          <div className="pt-4 space-y-3">
            <div className="text-xs font-mono uppercase tracking-wider text-zinc-400 font-semibold mb-2">
              Upcoming
            </div>
            {planned.map((item, i) => (
              <motion.div
                key={item}
                initial={{ opacity: 0, y: 6 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.25, delay: i * 0.02 }}
                className="flex items-start gap-3 text-sm text-zinc-400 leading-snug"
              >
                <Square className="w-4 h-4 text-zinc-300 shrink-0 mt-0.5" />
                <span>{item}</span>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
