import { Navbar } from "@/components/Navbar";
import { Hero } from "@/components/Hero";
import { Highlights } from "@/components/Highlights";
import { UsageShowcase } from "@/components/UsageShowcase";
import { FeaturesChecklist } from "@/components/FeaturesChecklist";
import { Footer } from "@/components/Footer";

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen bg-white">
      <Navbar />
      <main className="flex-1">
        <Hero />
        <Highlights />
        <UsageShowcase />
        <FeaturesChecklist />
      </main>
      <Footer />
    </div>
  );
}
