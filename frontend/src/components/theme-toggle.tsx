"use client";

import { MoonIcon, SunIcon } from "@heroicons/react/24/outline";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return (
      <div className="w-9 h-9 rounded-lg bg-muted animate-pulse" />
    );
  }

  const isDark = resolvedTheme === "dark";

  return (
    <button
      onClick={() => setTheme(isDark ? "light" : "dark")}
      className="relative w-9 h-9 rounded-lg bg-secondary hover:bg-secondary/80 transition-all duration-300 flex items-center justify-center group"
      aria-label="Toggle theme"
    >
      <div className="relative w-5 h-5">
        <SunIcon
          className={`size-5 text-amber-500 absolute inset-0 transition-all duration-300 ${
            isDark
              ? "opacity-0 rotate-90 scale-50"
              : "opacity-100 rotate-0 scale-100"
          }`}
        />
        <MoonIcon
          className={`size-5 text-violet-400 absolute inset-0 transition-all duration-300 ${
            isDark
              ? "opacity-100 rotate-0 scale-100"
              : "opacity-0 -rotate-90 scale-50"
          }`}
        />
      </div>
      <div
        className={`absolute inset-0 rounded-lg transition-all duration-300 ${
          isDark
            ? "bg-violet-500/10 ring-1 ring-violet-500/20"
            : "bg-amber-500/10 ring-1 ring-amber-500/20"
        }`}
      />
    </button>
  );
}