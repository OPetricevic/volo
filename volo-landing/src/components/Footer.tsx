export function Footer() {
  return (
    <footer className="border-t border-white/[0.06] py-10 md:py-12">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-2.5">
            <div className="w-6 h-6 bg-gradient-to-br from-sky-400 to-blue-600 rounded-md flex items-center justify-center">
              <span className="text-white text-[9px] font-bold">V</span>
            </div>
            <span className="text-sm font-semibold text-white/80">Volo</span>
            <span className="text-xs text-white/25 ml-2 hidden sm:inline">Voice-first browser assistant</span>
          </div>

          <div className="flex items-center gap-6 text-sm text-white/30">
            <a href="https://github.com/OPetricevic/volo" className="hover:text-white/70 transition-colors">GitHub</a>
            <a href="https://github.com/OPetricevic/volo#readme" className="hover:text-white/70 transition-colors">Docs</a>
          </div>
        </div>

        <div className="mt-8 text-center text-xs text-white/15">
          Bachelor Thesis · Omar Petricevic · 2025
        </div>
      </div>
    </footer>
  );
}
