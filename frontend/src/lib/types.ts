export interface JobMeta {
  location?: string;
  remote?: boolean;
  department?: string;
  team?: string;
  employmentType?: string;
  compensation?: string;
  jobUrl?: string;
  contentHash?: string;
  source: string;
  rawData?: unknown;
}

export interface Job {
  id: number;
  jobName: string;
  description: string;
  date: string;
  applyLink: string;
  companyName: string;
  compensation?: string;
  score?: number;
  meta: JobMeta;
}

export interface JobsResponse {
  jobs: Job[];
  offset: number;
  limit: number;
  total: number;
}

export interface SearchResult {
  id: string;
  score: number;
  job_id: string;
  company: string;
  title: string;
  location: string;
  team: string;
  department: string;
  employment_type: string;
  remote: boolean;
  description: string;
  apply_url: string;
  job_url: string;
  compensation: string;
  content_hash: string;
  is_active: boolean;
  status: string;
  published_at: string;
  scraped_at: string;
}

export interface SearchResponse {
  results: SearchResult[];
  query: string;
  total: number;
}