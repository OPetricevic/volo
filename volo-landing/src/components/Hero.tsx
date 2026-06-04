import { useState, useEffect } from "react";
import { HeroBackground } from "./HeroBackground";

export function Hero() {
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => setLoaded(true), 600);
    return () => clearTimeout(timer);
  }, []);

  return (
    <section className="pt-28 pb-16 md:pt-44 md:pb-32 relative overflow-hidden min-h-[90vh] md:min-h-screen flex items-center">
      {/* Three.js WebGL background — fades in after text loads */}
      <div
        className={`absolute inset-0 transition-opacity duration-[3000ms] ease-out ${loaded ? 'opacity-100' : 'opacity-0'}`}
      >
        <HeroBackground />
      </div>

      {/* Vignette on top of WebGL — strong center darkening for text readability */}
      <div className="absolute inset-0 pointer-events-none" style={{
        background: 'radial-gradient(ellipse 80% 70% at 50% 50%, rgba(6,6,9,0.7) 0%, rgba(6,6,9,0.4) 40%, rgba(6,6,9,0.85) 100%)'
      }} />

      {/* Top fade for nav readability */}
      <div className="absolute top-0 inset-x-0 h-40 bg-gradient-to-b from-[#060609] via-[#060609]/80 to-transparent pointer-events-none" />

      {/* Bottom fade */}
      <div className="absolute bottom-0 inset-x-0 h-32 bg-gradient-to-t from-[#060609] to-transparent pointer-events-none" />

      <div className="max-w-4xl mx-auto px-6 text-center relative z-10" style={{ textShadow: '0 2px 20px rgba(0,0,0,0.8)' }}>
        {/* Badge */}
        <div
          className={`inline-flex items-center gap-2 px-4 py-1.5 glass-card rounded-full text-xs text-white/50 mb-6 md:mb-8 transition-all duration-[1200ms] ease-out ${loaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-3'}`}
          style={{ transitionDelay: '0ms' }}
        >
          <span className="w-1.5 h-1.5 bg-emerald-400 rounded-full animate-pulse shadow-[0_0_6px_rgba(52,211,153,0.6)]" />
          Available for Chrome, Opera & Edge
        </div>

        {/* Headline — always visible immediately */}
        <h1
          className={`text-4xl sm:text-5xl md:text-7xl lg:text-8xl font-bold tracking-tight leading-[1.05] transition-all duration-[1200ms] ease-out ${loaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'}`}
          style={{ transitionDelay: '300ms' }}
        >
          <span className="text-white">Your browser,</span>
          <br />
          <span className="text-white">
            voice-first.
          </span>
        </h1>

        {/* Subtitle */}
        <p
          className={`mt-6 md:mt-8 text-base md:text-xl text-white max-w-2xl mx-auto leading-relaxed px-2 transition-all duration-[1200ms] ease-out ${loaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'}`}
          style={{ transitionDelay: '700ms' }}
        >
          A voice-powered browser assistant that searches, navigates, and learns your habits.
        </p>

        {/* CTAs */}
        <div
          className={`mt-10 md:mt-12 flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-4 transition-all duration-[1200ms] ease-out ${loaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'}`}
          style={{ transitionDelay: '1100ms' }}
        >
          <a
            href="#download"
            className="w-full sm:w-auto inline-flex items-center justify-center gap-2.5 px-7 py-3.5 bg-white text-[#060609] font-semibold rounded-full hover:shadow-[0_0_30px_rgba(255,255,255,0.15)] transition-all text-center text-[15px]"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <circle cx="8" cy="8" r="7" stroke="currentColor" strokeWidth="1.2"/>
              <circle cx="8" cy="8" r="3" fill="currentColor" opacity="0.9"/>
              <path d="M8 1a7 7 0 0 1 4.95 2.05L8 8" stroke="currentColor" strokeWidth="1.2" fill="none"/>
            </svg>
            Add to Browser
          </a>
          <a
            href="#download"
            className="w-full sm:w-auto inline-flex items-center justify-center gap-2.5 px-7 py-3.5 glass-card text-white/90 font-medium rounded-full hover:text-white hover:border-white/[0.15] transition-all text-center text-[15px]"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path d="M0 2.3l6.5-.9v6.3H0V2.3zm7.3-1L15.8.1v7.6H7.3V1.3zM15.8 8.4V16l-8.5-1.2V8.4h8.5zM6.5 14.7L0 13.8V8.4h6.5v6.3z"/>
            </svg>
            Download for Windows
          </a>
        </div>

        <p
          className={`mt-5 text-xs text-white/25 transition-all duration-[1200ms] ease-out ${loaded ? 'opacity-100' : 'opacity-0'}`}
          style={{ transitionDelay: '1500ms' }}
        >
          Free & open source · No account required
        </p>
      </div>
    </section>
  );
}
