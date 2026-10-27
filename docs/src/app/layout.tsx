import type { Metadata, Viewport } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
  display: "swap",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
  display: "swap",
});

const siteUrl = "https://jvm-cli.vercel.app";

export const viewport: Viewport = {
  themeColor: "#ffffff",
  width: "device-width",
  initialScale: 1,
};

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title: {
    default: "jvm — The Ultra-Fast Java Version Manager",
    template: "%s | jvm",
  },
  description:
    "An ultra-fast, cross-platform Java Version Manager written in Go with an nvm-style architecture. Sub-3ms shell startup time, zero dependencies, external JDK linking, and project auto-switching.",
  applicationName: "jvm",
  authors: [
    {
      name: "Tanmay",
      url: "https://github.com/unstoppable-tanmay",
    },
  ],
  generator: "Next.js",
  keywords: [
    "jvm",
    "jvm cli",
    "java version manager",
    "nvm for java",
    "jdk manager",
    "sdkman alternative",
    "openjdk manager",
    "adoptium temurin",
    "amazon corretto",
    "azul zulu",
    "developer tools",
    "go cli",
  ],
  referrer: "origin-when-cross-origin",
  creator: "Tanmay",
  publisher: "Tanmay",
  formatDetection: {
    email: false,
    address: false,
    telephone: false,
  },
  icons: {
    icon: [
      { url: "/favicon.svg", type: "image/svg+xml" },
    ],
    apple: "/favicon.svg",
  },
  openGraph: {
    title: "jvm — The Ultra-Fast Java Version Manager",
    description:
      "A modern, lightweight Java Version Manager written in Go with an nvm-style architecture. Sub-3ms startup, zero external dependencies, native Windows support.",
    url: siteUrl,
    siteName: "jvm",
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "jvm — The Ultra-Fast Java Version Manager",
    description:
      "A modern, lightweight Java Version Manager written in Go with an nvm-style architecture.",
    creator: "@tanmay",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-video-preview": -1,
      "max-image-preview": "large",
      "max-snippet": -1,
    },
  },
  alternates: {
    canonical: siteUrl,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable}`}>
      <head>
        <link rel="icon" href="/favicon.svg" type="image/svg+xml" />
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              "@context": "https://schema.org",
              "@type": "SoftwareApplication",
              name: "jvm",
              operatingSystem: "macOS, Linux, Windows",
              applicationCategory: "DeveloperApplication",
              description:
                "An ultra-fast Java Version Manager written in Go with an nvm-style architecture.",
              url: "https://github.com/unstoppable-tanmay/jvm",
              author: {
                "@type": "Person",
                name: "Tanmay",
                url: "https://github.com/unstoppable-tanmay",
              },
              offers: {
                "@type": "Offer",
                price: "0",
                priceCurrency: "USD",
              },
            }),
          }}
        />
      </head>
      <body className="bg-white text-zinc-900 antialiased selection:bg-zinc-900 selection:text-white min-h-screen">
        {children}
      </body>
    </html>
  );
}
