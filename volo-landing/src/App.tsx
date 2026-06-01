import { Nav } from "./components/Nav";
import { Hero } from "./components/Hero";
import { ProductShowcase } from "./components/ProductShowcase";
import { Features } from "./components/Features";
import { HowItWorks } from "./components/HowItWorks";
import { Download } from "./components/Download";
import { Footer } from "./components/Footer";

export function App() {
  return (
    <div className="min-h-screen">
      <Nav />
      <Hero />
      <ProductShowcase />
      <Features />
      <HowItWorks />
      <Download />
      <Footer />
    </div>
  );
}
