"use client";

import { useEffect, useState } from "react";
import { Job } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import {
  MapPinIcon,
  GlobeAltIcon,
  BuildingOfficeIcon,
  ArrowTopRightOnSquareIcon,
  XMarkIcon,
  CheckIcon,
  LinkIcon,
  CalendarDaysIcon,
  BookmarkIcon,
  CurrencyDollarIcon,
} from "@heroicons/react/24/outline";

function formatFullDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

export function JobSidebar({
  job,
  onClose,
}: {
  job: Job;
  onClose: () => void;
}) {
  const [copied, setCopied] = useState(false);

  const jobUrl = typeof window !== "undefined"
    ? `${window.location.origin}/?job=${job.id}`
    : `/?job=${job.id}`;

  const handleShare = async () => {
    try {
      await navigator.clipboard.writeText(jobUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy:", err);
    }
  };

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    document.addEventListener("keydown", handleEscape);
    return () => document.removeEventListener("keydown", handleEscape);
  }, [onClose]);

  useEffect(() => {
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = "";
    };
  }, []);

  return (
    <>
      <div
        className="fixed inset-0 z-40 bg-black/20 backdrop-blur-sm animate-fade-in"
        onClick={onClose}
      />

      <aside className="fixed right-0 top-0 z-50 h-full w-full sm:w-[60%] bg-background border-l shadow-2xl animate-slide-in overflow-y-auto">
        <div className="sticky top-0 z-10 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 border-b">
          <div className="flex items-center justify-between p-4">
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="text-[10px] px-2 py-0.5">
                {job.meta.source}
              </Badge>
              <span className="text-xs text-muted-foreground">
                {formatFullDate(job.date)}
              </span>
            </div>
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={onClose}
              className="rounded-full"
            >
              <XMarkIcon className="size-4" />
            </Button>
          </div>
        </div>

        <div className="p-6 space-y-6">
          <div>
            <h2 className="text-xl font-semibold leading-tight mb-2">
              {job.jobName}
            </h2>
            <p className="text-muted-foreground">{job.companyName}</p>
          </div>

          <div className="flex flex-wrap gap-2">
            {job.meta.location && (
              <Badge variant="secondary" className="text-xs gap-1.5 px-3 py-1">
                <MapPinIcon className="size-3.5" />
                {job.meta.location}
              </Badge>
            )}
            {job.meta.remote && (
              <Badge variant="secondary" className="text-xs gap-1.5 px-3 py-1">
                <GlobeAltIcon className="size-3.5" />
                Remote
              </Badge>
            )}
            {job.compensation && (
              <Badge variant="secondary" className="text-xs gap-1.5 px-3 py-1 bg-emerald-50 text-emerald-700 border-emerald-200">
                <CurrencyDollarIcon className="size-3.5" />
                {job.compensation}
              </Badge>
            )}
            {job.meta.employmentType && (
              <Badge variant="outline" className="text-xs px-3 py-1">
                {job.meta.employmentType}
              </Badge>
            )}
            {job.meta.department && (
              <Badge variant="outline" className="text-xs px-3 py-1">
                {job.meta.department}
              </Badge>
            )}
          </div>

          <Separator />

          <div className="space-y-4">
            <h3 className="text-sm font-semibold">About this role</h3>
            <p className="text-sm text-muted-foreground leading-relaxed whitespace-pre-wrap">
              {job.description}
            </p>
          </div>

          <Separator />

          <div className="space-y-3">
            <h3 className="text-sm font-semibold">Details</h3>
            <div className="grid grid-cols-2 gap-3">
              {job.meta.team && (
                <div className="flex items-start gap-2.5 text-sm">
                  <BuildingOfficeIcon className="size-4 text-muted-foreground mt-0.5" />
                  <div>
                    <p className="text-muted-foreground text-xs">Team</p>
                    <p className="font-medium">{job.meta.team}</p>
                  </div>
                </div>
              )}
              <div className="flex items-start gap-2.5 text-sm">
                <CalendarDaysIcon className="size-4 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-muted-foreground text-xs">Posted</p>
                  <p className="font-medium">{formatFullDate(job.date)}</p>
                </div>
              </div>
              <div className="flex items-start gap-2.5 text-sm">
                <BookmarkIcon className="size-4 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-muted-foreground text-xs">Source</p>
                  <p className="font-medium">{job.meta.source}</p>
                </div>
              </div>
              {job.meta.remote !== undefined && (
                <div className="flex items-start gap-2.5 text-sm">
                  <GlobeAltIcon className="size-4 text-muted-foreground mt-0.5" />
                  <div>
                    <p className="text-muted-foreground text-xs">Work Type</p>
                    <p className="font-medium">{job.meta.remote ? "Remote" : "On-site"}</p>
                  </div>
                </div>
              )}
            </div>
          </div>

          <Separator />

          <div className="flex gap-3">
            <Button asChild className="flex-1 h-11 rounded-xl">
              <a href={job.applyLink} target="_blank" rel="noopener noreferrer">
                Apply Now
                <ArrowTopRightOnSquareIcon className="size-4" />
              </a>
            </Button>
            <Button
              variant="outline"
              size="icon"
              onClick={handleShare}
              className="h-11 w-11 rounded-xl"
            >
              {copied ? (
                <CheckIcon className="size-4 text-green-500" />
              ) : (
                <LinkIcon className="size-4" />
              )}
            </Button>
          </div>

          {copied && (
            <p className="text-xs text-center text-green-500 animate-fade-in">
              Link copied to clipboard!
            </p>
          )}
        </div>
      </aside>
    </>
  );
}
