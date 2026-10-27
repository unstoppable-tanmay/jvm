export const HERO_CONTENT = {
  badge: "Sub-3ms Shell Startup",
  title: "Java version management,\nmade simple.",
  subtitle:
    "Stop wrestling with broken JAVA_HOME paths, slow shell startups, and manual downloads. Switch between Java 8, 17, 21, and beyond in milliseconds.",
  installCommands: {
    mac: "curl -fsSL https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.sh | bash",
    brew: "brew tap unstoppable-tanmay/jvm && brew install jvm",
    win: "irm https://raw.githubusercontent.com/unstoppable-tanmay/jvm/main/install.ps1 | iex",
  },
  features: [
    "⚡ Sub-3ms Shell Startup",
    "🔗 Link Existing JDKs",
    "📁 Auto-Switch with .java-version",
    "🪟 Native Windows Support",
  ],
};
