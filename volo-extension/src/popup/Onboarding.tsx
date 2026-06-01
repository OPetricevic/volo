import { useState } from "react";
import type { MicMode } from "@/shared/types";

interface OnboardingProps {
  onComplete: (micMode: MicMode) => void;
}

export function Onboarding({ onComplete }: OnboardingProps) {
  const [step, setStep] = useState<"welcome" | "mic">("welcome");

  return (
    <div className="w-[320px] bg-[#1e1f22] text-[#dbdee1] p-6">
      {step === "welcome" ? (
        <WelcomeStep onNext={() => setStep("mic")} />
      ) : (
        <MicStep onComplete={onComplete} />
      )}
    </div>
  );
}

function WelcomeStep({ onNext }: { onNext: () => void }) {
  return (
    <div className="text-center">
      <div className="w-14 h-14 bg-[#5865f2] rounded-2xl flex items-center justify-center mx-auto mb-5">
        <span className="text-white text-xl font-bold">V</span>
      </div>
      <h1 className="text-[18px] font-semibold mb-2">Welcome to Volo</h1>
      <p className="text-[13px] text-[#949ba4] leading-relaxed mb-6">
        Your voice-first browser assistant. Say "Hey Volo" and tell it what you want — search, navigate, open apps.
      </p>
      <button
        onClick={onNext}
        className="w-full py-2.5 bg-[#5865f2] text-white text-[13px] font-medium rounded-lg hover:bg-[#4752c4] transition"
      >
        Get Started
      </button>
    </div>
  );
}

function MicStep({ onComplete }: { onComplete: (mode: MicMode) => void }) {
  return (
    <div>
      <div className="text-center mb-5">
        <div className="w-12 h-12 bg-[#2b2d31] rounded-full flex items-center justify-center mx-auto mb-4">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#5865f2" strokeWidth="1.5" strokeLinecap="round">
            <rect x="8" y="3" width="8" height="12" rx="4"/>
            <path d="M5 11a7 7 0 0 0 14 0"/>
            <line x1="12" y1="18" x2="12" y2="21"/>
          </svg>
        </div>
        <h2 className="text-[16px] font-semibold mb-1">Microphone Access</h2>
        <p className="text-[12px] text-[#949ba4] leading-relaxed">
          Volo needs your microphone to listen for voice commands. Choose how you'd like it to work:
        </p>
      </div>

      <div className="space-y-2 mb-5">
        <MicChoice
          title="Always Listen"
          description="Continuously listens for 'Hey Volo'. Best experience."
          recommended
          onClick={() => onComplete("always")}
        />
        <MicChoice
          title="Push to Talk"
          description="Only listens when you click the icon or press Ctrl+Shift+V."
          onClick={() => onComplete("once")}
        />
        <MicChoice
          title="Keep Off"
          description="No microphone access. You can enable it later in settings."
          onClick={() => onComplete("off")}
        />
      </div>

      <p className="text-[10px] text-[#6d6f78] text-center">
        You can change this anytime in Volo settings.
      </p>
    </div>
  );
}

function MicChoice({
  title,
  description,
  recommended,
  onClick,
}: {
  title: string;
  description: string;
  recommended?: boolean;
  onClick: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className={`w-full text-left p-3 rounded-lg border transition hover:border-[#5865f2]/50 hover:bg-[#5865f2]/5 ${
        recommended ? "border-[#5865f2]/30 bg-[#5865f2]/5" : "border-[#3f4147]/50"
      }`}
    >
      <div className="flex items-center gap-2">
        <span className="text-[13px] font-medium text-[#dbdee1]">{title}</span>
        {recommended && (
          <span className="text-[9px] px-1.5 py-0.5 bg-[#5865f2]/20 text-[#5865f2] rounded font-medium">Recommended</span>
        )}
      </div>
      <p className="text-[11px] text-[#949ba4] mt-0.5">{description}</p>
    </button>
  );
}
