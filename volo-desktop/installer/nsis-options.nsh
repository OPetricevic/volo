; ─────────────────────────────────────────────────────────────
; Volo Desktop — NSIS Installer Custom Options
; Downloads Ollama during install, then pulls Gemma 2B model.
; Uses inetc plugin (bundled with electron-builder's NSIS).
; ─────────────────────────────────────────────────────────────

!macro customInstall
  ; ── Step 1: Check if Ollama is already installed ──
  IfFileExists "$LOCALAPPDATA\Programs\Ollama\ollama.exe" ollama_exists 0

  ; ── Step 2: Download Ollama installer ──
  DetailPrint "Downloading Volo AI engine (Ollama)..."
  DetailPrint "This requires an internet connection."

  inetc::get "https://ollama.com/download/OllamaSetup.exe" "$TEMP\OllamaSetup.exe" /END
  Pop $0
  StrCmp $0 "OK" +3
    DetailPrint "Download failed ($0). AI features will need manual Ollama install."
    Goto done

  ; ── Step 3: Run Ollama installer silently ──
  DetailPrint "Installing Ollama..."
  ExecWait '"$TEMP\OllamaSetup.exe" /VERYSILENT /NORESTART /SUPPRESSMSGBOXES' $0
  Delete "$TEMP\OllamaSetup.exe"

  IntCmp $0 0 +2
    DetailPrint "Ollama install returned code $0. Continuing..."

  ollama_exists:
  DetailPrint "Ollama is installed."

  ; ── Step 4: Wait for Ollama to be ready ──
  DetailPrint "Starting Ollama..."
  Sleep 3000

  ; ── Step 5: Pull Gemma 2B model ──
  DetailPrint "Downloading AI model (Gemma 2B ~1.6GB)..."
  DetailPrint "This will take a few minutes. Please wait."
  nsExec::ExecToLog '"$LOCALAPPDATA\Programs\Ollama\ollama.exe" pull gemma2:2b'
  Pop $0

  StrCmp $0 "0" +3
    DetailPrint "Model download returned code $0. Volo will retry on first launch."
    Goto done

  DetailPrint "AI model ready."

  done:
!macroend

!macro customUnInstall
  ; Don't remove Ollama — user may use it independently.
  MessageBox MB_OK "Note: Ollama and AI models are not removed. Uninstall Ollama separately from Settings > Apps if desired."
!macroend
