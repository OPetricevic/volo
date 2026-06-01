export function Hero() {
  return (
    <section className="pt-32 pb-20 md:pt-44 md:pb-32 relative overflow-hidden min-h-screen flex items-center">
      {/* Raycast-style animated background rays */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        {/* Main glow */}
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[900px] h-[900px]">
          {/* Ray streaks */}
          <div className="absolute inset-0 animate-slow-spin">
            <div className="absolute top-[10%] left-[20%] w-[60%] h-[8%] bg-gradient-to-r from-transparent via-sky-500/30 to-transparent rotate-[35deg] blur-[30px]" />
            <div className="absolute top-[25%] left-[15%] w-[70%] h-[6%] bg-gradient-to-r from-transparent via-cyan-400/20 to-transparent rotate-[40deg] blur-[25px]" />
            <div className="absolute top-[40%] left-[10%] w-[80%] h-[10%] bg-gradient-to-r from-transparent via-sky-400/25 to-transparent rotate-[32deg] blur-[35px]" />
            <div className="absolute top-[55%] left-[18%] w-[65%] h-[7%] bg-gradient-to-r from-transparent via-teal-400/20 to-transparent rotate-[38deg] blur-[28px]" />
            <div className="absolute top-[70%] left-[12%] w-[75%] h-[9%] bg-gradient-to-r from-transparent via-sky-300/15 to-transparent rotate-[42deg] blur-[32px]" />
          </div>

          {/* Orb glow */}
          <div className="absolute top-[30%] left-[20%] w-[200px] h-[200px] bg-sky-500/20 rounded-full blur-[80px] animate-pulse-slow" />
          <div className="absolute bottom-[30%] right-[25%] w-[150px] h-[150px] bg-teal-500/15 rounded-full blur-[60px] animate-pulse-slow" style={{ animationDelay: '2s' }} />
        </div>

        {/* Vignette overlay */}
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_center,transparent_30%,#0a0a0f_75%)]" />
      </div>

      <div className="max-w-4xl mx-auto px-6 text-center relative z-10">
        {/* Badge */}
        <div className="inline-flex items-center gap-2 px-3 py-1 bg-volo-surface/80 border border-volo-border rounded-full text-xs text-volo-muted mb-8 backdrop-blur-sm">
          <span className="w-1.5 h-1.5 bg-volo-accent rounded-full animate-pulse" />
          Now available for Chrome, Opera & Edge
        </div>

        {/* Headline */}
        <h1 className="text-5xl md:text-7xl lg:text-8xl font-bold tracking-tight leading-[1.0]">
          Your browser,
          <br />
          <span className="bg-gradient-to-r from-volo-accent via-cyan-300 to-volo-teal bg-clip-text text-transparent">
            voice-first.
          </span>
        </h1>

        {/* Subtitle */}
        <p className="mt-8 text-lg md:text-xl text-volo-muted max-w-2xl mx-auto leading-relaxed">
          Say <span className="text-volo-text font-medium">"Hey Volo"</span> and tell it what you want.
          Search the web, open apps, navigate sites. It learns your patterns and gets faster over time.
        </p>

        {/* CTAs */}
        <div className="mt-12 flex flex-col sm:flex-row items-center justify-center gap-4">
          <a
            href="#download"
            className="px-8 py-3.5 bg-white text-black font-semibold rounded-full hover:bg-gray-100 transition-all shadow-lg shadow-white/10"
          >
            Download for Windows
          </a>
          <a
            href="#download"
            className="px-8 py-3.5 bg-volo-surface/80 border border-volo-border text-volo-text font-medium rounded-full hover:border-volo-muted backdrop-blur-sm transition-all"
          >
            Add to Chrome
          </a>
        </div>

        <p className="mt-5 text-xs text-volo-muted">
          v0.1.0 · Free & open source · No account required
        </p>
      </div>
    </section>
  );
}
