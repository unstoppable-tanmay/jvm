import { ImageResponse } from "next/og";
import { HERO_CONTENT } from "@/content/hero";

export async function GET() {
  let version = "v0.1.0";
  try {
    const res = await fetch(
      "https://api.github.com/repos/unstoppable-tanmay/jvm/releases/latest",
      { headers: { "User-Agent": "jvm-docs" } }
    );
    if (res.ok) {
      const data = await res.json();
      if (data?.tag_name) {
        version = data.tag_name;
      }
    }
  } catch {
    // fallback
  }

  return new ImageResponse(
    (
      <div
        style={{
          width: "100%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
          padding: "70px 80px",
          backgroundColor: "#ffffff",
        }}
      >
        {/* Top Header */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            width: "100%",
          }}
        >
          <div
            style={{
              display: "flex",
              fontSize: 40,
              fontWeight: 700,
              color: "#09090b",
              letterSpacing: "-0.04em",
            }}
          >
            jvm
          </div>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "12px",
              fontSize: 18,
              color: "#52525b",
              padding: "8px 20px",
              borderRadius: "9999px",
              backgroundColor: "#f4f4f5",
              fontWeight: 500,
            }}
          >
            <span>{version}</span>
            <span>•</span>
            <span>{HERO_CONTENT.badge}</span>
          </div>
        </div>

        {/* Center Hero Content */}
        <div style={{ display: "flex", flexDirection: "column", gap: "18px" }}>
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              fontSize: 64,
              fontWeight: 700,
              color: "#09090b",
              lineHeight: 1.1,
              letterSpacing: "-0.04em",
            }}
          >
            <span>Java version management,</span>
            <span>made simple.</span>
          </div>
          <div
            style={{
              display: "flex",
              fontSize: 24,
              color: "#52525b",
              maxWidth: "960px",
              lineHeight: 1.45,
            }}
          >
            {HERO_CONTENT.subtitle}
          </div>
        </div>

        {/* Dynamic feature highlights */}
        <div
          style={{
            display: "flex",
            gap: "28px",
            borderTop: "1px solid #f4f4f5",
            paddingTop: "28px",
            width: "100%",
          }}
        >
          {HERO_CONTENT.features.map((feature) => (
            <div
              key={feature}
              style={{
                display: "flex",
                alignItems: "center",
                fontSize: 18,
                color: "#27272a",
                fontWeight: 500,
              }}
            >
              <span>{feature}</span>
            </div>
          ))}
        </div>
      </div>
    ),
    {
      width: 1200,
      height: 630,
    }
  );
}
