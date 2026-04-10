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
  meta: JobMeta;
}

export interface JobsResponse {
  jobs: Job[];
  offset: number;
  limit: number;
  total: number;
}