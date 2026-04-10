"use client";

import { useState, useEffect, useMemo, useCallback, useRef } from "react";
import { useSearchParams } from "next/navigation";
import { Job, JobsResponse } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";
import { Separator } from "@/components/ui/separator";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Empty,
  EmptyDescription,
} from "@/components/ui/empty";
import { JobSidebar } from "@/components/job-sidebar";
import {
  MagnifyingGlassIcon,
  MapPinIcon,
  BriefcaseIcon,
  GlobeAltIcon,
  BuildingOfficeIcon,
  ArrowTopRightOnSquareIcon,
  ArrowPathIcon,
  XMarkIcon,
  HeartIcon,
  EyeIcon,
} from "@heroicons/react/24/outline";
import { ThemeToggle } from "@/components/theme-toggle";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const PAGE_SIZE = 80;

function formatRelativeDate(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffTime = now.getTime() - date.getTime();
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

  if (diffDays === 0) return "Today";
  if (diffDays === 1) return "Yesterday";
  if (diffDays < 7) return `${diffDays}d ago`;
  if (diffDays < 30) return `${Math.floor(diffDays / 7)}w ago`;
  if (diffDays < 365) return `${Math.floor(diffDays / 30)}mo ago`;
  return `${Math.floor(diffDays / 365)}y ago`;
}

function getInitials(companyName: string): string {
  return companyName
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((n) => n[0])
    .join("")
    .toUpperCase();
}

function JobCard({
  job,
  className = "",
  onView,
}: {
  job: Job;
  className?: string;
  onView: (job: Job) => void;
}) {
  const initials = getInitials(job.companyName);

  return (
    <Card className={`job-card border-muted/50 ${className}`}>
      <CardHeader className="flex flex-row items-start gap-3 pb-3">
        <Avatar className="size-10 rounded-lg">
          <AvatarFallback className="bg-muted text-sm font-semibold">
            {initials}
          </AvatarFallback>
        </Avatar>
        <div className="flex-1 min-w-0">
          <CardTitle className="text-sm font-semibold leading-snug line-clamp-2">
            {job.jobName}
          </CardTitle>
          <p className="text-xs text-muted-foreground mt-0.5 truncate">
            {job.companyName}
          </p>
        </div>
      </CardHeader>
      <CardContent className="space-y-3 pt-0">
        <p className="text-xs text-muted-foreground line-clamp-2 leading-relaxed">
          {job.description}
        </p>

        <div className="flex flex-wrap gap-1">
          {job.meta.location && (
            <Badge variant="secondary" className="text-[10px] px-1.5 py-0.5 gap-1">
              <MapPinIcon className="size-2.5" />
              {job.meta.location}
            </Badge>
          )}
          {job.meta.remote && (
            <Badge variant="secondary" className="text-[10px] px-1.5 py-0.5 gap-1">
              <GlobeAltIcon className="size-2.5" />
              Remote
            </Badge>
          )}
          {job.meta.employmentType && (
            <Badge variant="outline" className="text-[10px] px-1.5 py-0.5">
              {job.meta.employmentType}
            </Badge>
          )}
          {job.meta.department && (
            <Badge variant="outline" className="text-[10px] px-1.5 py-0.5">
              {job.meta.department}
            </Badge>
          )}
        </div>

        <Separator />

        <div className="flex items-center justify-between gap-2">
          <div className="flex items-center gap-2 text-[10px] text-muted-foreground min-w-0">
            <span className="flex items-center gap-1 truncate">
              <BuildingOfficeIcon className="size-2.5 shrink-0" />
              {job.meta.team || job.meta.source}
            </span>
            <span className="text-muted-foreground/60">•</span>
            <span className="shrink-0">{formatRelativeDate(job.date)}</span>
          </div>
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => onView(job)}
              className="h-7 w-7 rounded-lg"
            >
              <EyeIcon className="size-3.5" />
            </Button>
            <Button asChild size="sm" className="h-7 text-xs gap-1">
              <a href={job.applyLink} target="_blank" rel="noopener noreferrer">
                Apply
                <ArrowTopRightOnSquareIcon className="size-2.5" />
              </a>
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

function JobSkeleton() {
  return (
    <Card>
      <CardHeader className="flex flex-row items-start gap-3 pb-3">
        <Skeleton className="size-10 rounded-lg" />
        <div className="flex-1 space-y-2">
          <Skeleton className="h-3.5 w-3/4" />
          <Skeleton className="h-2.5 w-1/2" />
        </div>
      </CardHeader>
      <CardContent className="space-y-3 pt-0">
        <Skeleton className="h-2.5 w-full" />
        <Skeleton className="h-2.5 w-2/3" />
        <div className="flex gap-1.5">
          <Skeleton className="h-4 w-14" />
          <Skeleton className="h-4 w-12" />
        </div>
        <Separator />
        <div className="flex justify-between items-center">
          <Skeleton className="h-2.5 w-32" />
          <Skeleton className="h-7 w-14" />
        </div>
      </CardContent>
    </Card>
  );
}

function SearchBar({
  search,
  onSearchChange,
  onSearch,
  selectedLocation,
  onLocationChange,
  locations,
  selectedCompany,
  onCompanyChange,
  companies,
  selectedSort,
  onSortChange,
}: {
  search: string;
  onSearchChange: (v: string) => void;
  onSearch: () => void;
  selectedLocation: string;
  onLocationChange: (v: string) => void;
  locations: string[];
  selectedCompany: string;
  onCompanyChange: (v: string) => void;
  companies: string[];
  selectedSort: string;
  onSortChange: (v: string) => void;
}) {
  return (
    <div className="flex flex-col sm:flex-row gap-2">
      <div className="relative flex-1">
        <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground pointer-events-none" />
        <Input
          placeholder="Search jobs, companies..."
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && onSearch()}
          className="pl-9 pr-16 h-10 bg-muted/30 border-0 focus:bg-background rounded-xl"
        />
        {search && (
          <button
            onClick={() => onSearchChange("")}
            className="absolute right-10 top-1/2 -translate-y-1/2 p-1 hover:bg-muted rounded-full transition-colors"
          >
            <XMarkIcon className="size-3 text-muted-foreground" />
          </button>
        )}
        <Button
          size="sm"
          onClick={onSearch}
          className="absolute right-1 top-1/2 -translate-y-1/2 h-7 px-3 rounded-lg"
        >
          Search
        </Button>
      </div>

      <div className="flex gap-2">
        <Select value={selectedCompany} onValueChange={onCompanyChange}>
          <SelectTrigger className="w-full sm:w-44 h-10 bg-muted/30 border-0 rounded-xl">
            <SelectValue placeholder="Company" />
          </SelectTrigger>
          <SelectContent position="popper" className="w-48">
            <SelectItem value="all">All Companies</SelectItem>
            {companies.map((company) => (
              <SelectItem key={company} value={company}>
                {company}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Select value={selectedLocation} onValueChange={onLocationChange}>
          <SelectTrigger className="w-full sm:w-40 h-10 bg-muted/30 border-0 rounded-xl">
            <MapPinIcon className="size-4 mr-2 text-muted-foreground" />
            <SelectValue placeholder="Location" />
          </SelectTrigger>
          <SelectContent position="popper" className="w-48">
            <SelectItem value="all">All Locations</SelectItem>
            {locations.map((loc) => (
              <SelectItem key={loc} value={loc}>
                {loc}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Select value={selectedSort} onValueChange={onSortChange}>
          <SelectTrigger className="w-full sm:w-36 h-10 bg-muted/30 border-0 rounded-xl">
            <SelectValue placeholder="Sort" />
          </SelectTrigger>
          <SelectContent position="popper">
            <SelectItem value="newest">Newest First</SelectItem>
            <SelectItem value="oldest">Oldest First</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}

export default function Page() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedLocation, setSelectedLocation] = useState("all");
  const [selectedCompany, setSelectedCompany] = useState("all");
  const [selectedSort, setSelectedSort] = useState("newest");
  const [locations, setLocations] = useState<string[]>([]);
  const [companies, setCompanies] = useState<string[]>([]);
  const [total, setTotal] = useState(0);
  const [offset, setOffset] = useState(0);
  const [selectedJob, setSelectedJob] = useState<Job | null>(null);
  const [pendingJobId, setPendingJobId] = useState<string | null>(null);
  const [isSearchingForJob, setIsSearchingForJob] = useState(false);
  const [invalidJobId, setInvalidJobId] = useState<string | null>(null);

  const searchParams = useSearchParams();
  const isInitialized = useRef(false);

  const fetchJobs = useCallback(
    async (pageOffset: number) => {
      setLoading(true);
      setError(null);
      try {
        const params = new URLSearchParams();
        params.set("offset", String(pageOffset));
        params.set("limit", String(PAGE_SIZE));
        if (searchQuery) params.set("search", searchQuery);
        if (selectedCompany !== "all") params.set("company", selectedCompany);
        if (selectedLocation !== "all") params.set("location", selectedLocation);
        params.set("sort", selectedSort);

        const res = await fetch(`${API_URL}/getallJobsFromSQL?${params.toString()}`);
        if (!res.ok) throw new Error("Failed to fetch jobs");
        const data: JobsResponse = await res.json();
        setJobs(data.jobs ?? []);
        setTotal(data.total);
        return data;
      } catch (e) {
        setError(e instanceof Error ? e.message : "Something went wrong");
        return null;
      } finally {
        setLoading(false);
      }
    },
    [searchQuery, selectedCompany, selectedLocation, selectedSort]
  );

  useEffect(() => {
    async function fetchFilters() {
      try {
        const [locationsRes, companiesRes] = await Promise.all([
          fetch(`${API_URL}/locations`),
          fetch(`${API_URL}/companies`),
        ]);
        if (locationsRes.ok) {
          const data = await locationsRes.json();
          setLocations(data.locations ?? []);
        }
        if (companiesRes.ok) {
          const data = await companiesRes.json();
          setCompanies(data.companies ?? []);
        }
      } catch {}
    }
    fetchFilters();
  }, []);

  useEffect(() => {
    fetchJobs(offset);
  }, [fetchJobs, offset]);

  const findJobById = useCallback(
    async (jobId: string): Promise<Job | null> => {
      setIsSearchingForJob(true);
      try {
        const res = await fetch(`${API_URL}/job/${jobId}`);
        if (!res.ok) {
          setIsSearchingForJob(false);
          return null;
        }
        const job: Job = await res.json();
        setIsSearchingForJob(false);
        return job;
      } catch {
        setIsSearchingForJob(false);
        return null;
      }
    },
    []
  );

  useEffect(() => {
    if (!isInitialized.current) return;
    const jobId = searchParams.get("job");
    if (jobId && !selectedJob && !isSearchingForJob) {
      const jobInList = jobs.find((j) => String(j.id) === jobId);
      if (jobInList) {
        setSelectedJob(jobInList);
      } else {
        setPendingJobId(jobId);
      }
    }
  }, [searchParams, jobs, selectedJob, isSearchingForJob]);

  useEffect(() => {
    if (!isInitialized.current) {
      isInitialized.current = true;
      const jobId = searchParams.get("job");
      if (jobId) {
        setPendingJobId(jobId);
      }
    }
  }, [searchParams]);

  useEffect(() => {
    if (pendingJobId && !isSearchingForJob) {
      findJobById(pendingJobId).then((foundJob) => {
        if (foundJob) {
          setSelectedJob(foundJob);
          setInvalidJobId(null);
        } else {
          setInvalidJobId(pendingJobId);
          const url = new URL(window.location.href);
          url.searchParams.delete("job");
          window.history.pushState({}, "", url);
        }
        setPendingJobId(null);
      });
    }
  }, [pendingJobId, isSearchingForJob, findJobById]);

  const handleLocationChange = (location: string) => {
    setSelectedLocation(location);
    setOffset(0);
  };

  const handleCompanyChange = (company: string) => {
    setSelectedCompany(company);
    setOffset(0);
  };

  const handleSortChange = (sort: string) => {
    setSelectedSort(sort);
    setOffset(0);
  };

  const handleSearch = () => {
    setOffset(0);
    setSearchQuery(search);
  };

  const handleViewJob = (job: Job) => {
    setSelectedJob(job);
    const url = new URL(window.location.href);
    url.searchParams.set("job", String(job.id));
    window.history.pushState({}, "", url);
  };

  const handleCloseSidebar = () => {
    setSelectedJob(null);
    const url = new URL(window.location.href);
    url.searchParams.delete("job");
    window.history.pushState({}, "", url);
  };

  const totalPages = Math.ceil(total / PAGE_SIZE);

  const filteredJobs = useMemo(() => {
    return jobs;
  }, [jobs]);

  const hasActiveFilters = searchQuery || selectedLocation !== "all" || selectedCompany !== "all";

  return (
    <div className="min-h-svh bg-background flex flex-col">
      <header className="sticky top-0 z-20 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/10 rounded-xl">
                <BriefcaseIcon className="size-5 text-primary" />
              </div>
              <div>
                <h1 className="text-xl font-semibold">Jobs</h1>
                <p className="text-xs text-muted-foreground">
                  {loading || isSearchingForJob
                    ? "Loading..."
                    : `${total} positions`}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <a
                href="https://github.com/MabudAlam/JobsScraper"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors px-2 py-1.5 rounded-lg hover:bg-secondary"
              >
                <svg className="size-4" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
                </svg>
                <span className="hidden sm:inline">JobsScraper</span>
              </a>
              <ThemeToggle />
            </div>
          </div>

          <SearchBar
            search={search}
            onSearchChange={setSearch}
            onSearch={handleSearch}
            selectedLocation={selectedLocation}
            onLocationChange={handleLocationChange}
            locations={locations}
            selectedCompany={selectedCompany}
            onCompanyChange={handleCompanyChange}
            companies={companies}
            selectedSort={selectedSort}
            onSortChange={handleSortChange}
          />
        </div>
      </header>

      <main className="flex-1 container mx-auto px-4 py-4">
        {error && (
          <Card className="border-destructive/50 bg-destructive/5 mb-4">
            <CardContent className="flex items-center justify-between gap-2 p-3 text-sm text-destructive">
              <span>{error}</span>
              <Button variant="ghost" size="sm" onClick={() => fetchJobs(offset)}>
                Retry
              </Button>
            </CardContent>
          </Card>
        )}

        {invalidJobId && (
          <Card className="border-amber-500/50 bg-amber-500/10 mb-4">
            <CardContent className="flex items-center justify-between gap-2 p-3 text-sm text-amber-700 dark:text-amber-400">
              <span>Job not found. It may have been removed or the link is invalid.</span>
              <Button variant="ghost" size="sm" onClick={() => setInvalidJobId(null)}>
                Dismiss
              </Button>
            </CardContent>
          </Card>
        )}

        {isSearchingForJob && (
          <Card className="mb-4">
            <CardContent className="flex items-center justify-center gap-2 p-4 text-sm text-muted-foreground">
              <ArrowPathIcon className="size-4 animate-spin" />
              <span>Finding job...</span>
            </CardContent>
          </Card>
        )}

        {loading ? (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <JobSkeleton key={i} />
            ))}
          </div>
        ) : filteredJobs.length === 0 ? (
          <Empty className="min-h-[250px]">
            <EmptyDescription>
              {hasActiveFilters
                ? "No jobs match your search."
                : "No jobs available at the moment."}
            </EmptyDescription>
          </Empty>
        ) : (
          <>
            <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
              {filteredJobs.map((job) => (
                <JobCard
                  key={job.id}
                  job={job}
                  className="job-card-enter"
                  onView={handleViewJob}
                />
              ))}
            </div>

            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-2 mt-6 pb-4">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setOffset((o) => Math.max(0, o - PAGE_SIZE))}
                  disabled={offset === 0 || loading || isSearchingForJob}
                  className="rounded-xl"
                >
                  <svg
                    className="size-3.5 mr-1"
                    viewBox="0 0 16 16"
                    fill="none"
                  >
                    <path
                      d="M10 12L6 8l4-4"
                      stroke="currentColor"
                      strokeWidth="1.5"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                  Prev
                </Button>
                <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                  <span>
                    Showing <span className="font-medium text-foreground">{offset + 1}</span>
                    {" - "}
                    <span className="font-medium text-foreground">{Math.min(offset + PAGE_SIZE, total)}</span>
                    {" of "}
                    <span className="font-medium text-foreground">{total}</span>
                    {" jobs"}
                  </span>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setOffset((o) => o + PAGE_SIZE)}
                  disabled={
                    offset + PAGE_SIZE >= total || loading || isSearchingForJob
                  }
                  className="rounded-xl"
                >
                  Next
                  <svg
                    className="size-3.5 ml-1"
                    viewBox="0 0 16 16"
                    fill="none"
                  >
                    <path
                      d="M6 4l4 4-4 4"
                      stroke="currentColor"
                      strokeWidth="1.5"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                </Button>
              </div>
            )}
          </>
        )}
      </main>

      <footer className="border-t bg-background/50">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-center relative">
            <p className="flex items-center gap-1.5 text-xs text-muted-foreground">
              Built with
              <HeartIcon className="size-3 text-red-500" />
              by Developers
            </p>
            <a
              href="https://github.com/MabudAlam/JobsScraper"
              target="_blank"
              rel="noopener noreferrer"
              className="absolute right-0 flex items-center gap-2 text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              <svg className="size-4" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z" />
              </svg>
              <span>JobsScraper</span>
            </a>
          </div>
        </div>
      </footer>

      {selectedJob && (
        <JobSidebar job={selectedJob} onClose={handleCloseSidebar} />
      )}
    </div>
  );
}
