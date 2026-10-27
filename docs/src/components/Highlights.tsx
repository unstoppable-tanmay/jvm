"use client";

import { motion } from "motion/react";

export function Highlights() {
  const points = [
    {
      title: "No More Shell Lag",
      description:
        "Tired of waiting for your terminal prompt every time you open a tab? jvm keeps shell startup under 3ms so your workflow stays snappy.",
    },
    {
      title: "Keep Your Existing JDKs",
      description:
        "Already installed Java via Homebrew, Android Studio, or macOS installers? Link them in seconds without redownloading 300MB files.",
    },
    {
      title: "Zero Setup Project Switching",
      description:
        "Commit a .java-version or .jvmrc file to your repository. When you or your teammates cd in, the right Java runtime activates automatically.",
    },
    {
      title: "Install from Any Source",
      description:
        "Fetch official builds from Temurin, Corretto, and Zulu in one command, or install internal corporate mirrors and early-access builds with checksum verification.",
    },
  ];

  return (
    <section id="highlights" className="py-14 border-y border-zinc-100 bg-zinc-50/50">
      <div className="max-w-4xl mx-auto px-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-8">
          {points.map((p, i) => (
            <motion.div
              key={p.title}
              initial={{ opacity: 0, y: 10 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.35, delay: i * 0.05 }}
              className="space-y-1.5"
            >
              <h3 className="text-base font-medium text-zinc-900 tracking-tight">
                {p.title}
              </h3>
              <p className="text-sm text-zinc-600 leading-relaxed">
                {p.description}
              </p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
