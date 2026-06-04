import { Nav } from "./components/Nav";
import { Hero } from "./components/Hero";
import { ProductShowcase } from "./components/ProductShowcase";
import { Features } from "./components/Features";
import { HowItWorks } from "./components/HowItWorks";
import { Download } from "./components/Download";
import { Footer } from "./components/Footer";

export function App() {
  return (
    <div className="min-h-screen relative overflow-hidden">
      {/* Noise overlay on the whole page */}
      <div className="fixed inset-0 pointer-events-none opacity-[0.015] transform-gpu" style={{
        backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 512 512' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.75' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E")`,
      }} />

      {/* Content */}
      <div className="relative z-10">
        <Nav />
        <Hero />

        {/* Ambient glow between Hero and Product Showcase */}
        <div className="relative pointer-events-none" aria-hidden="true">
          <div className="absolute -top-[200px] left-[5%] w-[500px] h-[500px] bg-sky-500/[0.03] rounded-full blur-[150px] transform-gpu" />
          <div className="absolute -top-[100px] right-[15%] w-[350px] h-[350px] bg-indigo-500/[0.025] rounded-full blur-[120px] transform-gpu" />
        </div>

        <ProductShowcase />

        {/* Ambient glow between Showcase and Features */}
        <div className="relative pointer-events-none" aria-hidden="true">
          <div className="absolute -top-[150px] right-[5%] w-[450px] h-[450px] bg-teal-500/[0.025] rounded-full blur-[130px] transform-gpu" />
          <div className="absolute -top-[250px] left-[20%] w-[300px] h-[300px] bg-cyan-500/[0.02] rounded-full blur-[100px] transform-gpu" />
        </div>

        <Features />

        {/* Ambient glow between Features and How It Works */}
        <div className="relative pointer-events-none" aria-hidden="true">
          <div className="absolute -top-[200px] left-[10%] w-[400px] h-[400px] bg-sky-400/[0.02] rounded-full blur-[120px] transform-gpu" />
          <div className="absolute -top-[100px] right-[25%] w-[350px] h-[350px] bg-violet-500/[0.02] rounded-full blur-[110px] transform-gpu" />
        </div>

        <HowItWorks />

        {/* Ambient glow before Download — slightly stronger to draw attention */}
        <div className="relative pointer-events-none" aria-hidden="true">
          <div className="absolute -top-[180px] left-[30%] w-[500px] h-[400px] bg-sky-500/[0.035] rounded-full blur-[140px] transform-gpu" />
          <div className="absolute -top-[100px] right-[10%] w-[300px] h-[300px] bg-teal-400/[0.025] rounded-full blur-[100px] transform-gpu" />
        </div>

        <Download />
        <Footer />
      </div>
    </div>
  );
}
