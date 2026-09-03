// Typed client for Go Fiber API — mirrors Pydantic/TS interfaces from schema

const BASE = import.meta.env.PUBLIC_API_URL || "http://localhost:8080";

export interface Person {
  id: string;
  org_id: string;
  first_name: string;
  middle_name?: string;
  last_name: string;
  preferred_name?: string;
  date_of_birth?: string;
  gender?: "male" | "female" | "non-binary" | "prefer-not-to-say";
  personal_email?: string;
  org_email?: string;
  phone_primary?: string;
  profile_photo_url?: string;
  city?: string;
  country?: string;
  is_international: boolean;
  is_active: boolean;
  tags: string[];
  current_job_title?: string;
  current_department?: string;
  current_team?: string;
  current_office_location?: string;
  current_hire_date?: string;
  created_at: string;
  updated_at: string;
}

export interface EmergencyContact {
  id: string;
  person_id: string;
  name: string;
  phone?: string;
  email?: string;
  relationship?: string;
  created_at: string;
  updated_at: string;
}

export interface EmploymentRecord {
  id: string;
  person_id: string;
  job_title?: string;
  job_level?: string;
  employment_status?: string;
  employment_type?: string;
  work_arrangement?: string;
  department?: string;
  team?: string;
  office_location?: string;
  salary_amount?: number;
  hire_date?: string;
  valid_from: string;
  valid_to?: string;
  is_current: boolean;
}



export interface StatusChangeEvent {
  id: string;
  person_id: string;
  org_id: string;
  event_type: string;
  context: "employment" | "general";
  from_department?: string;
  to_department?: string;
  reason?: string;
  is_voluntary?: boolean;
  effective_date: string;
  recorded_at: string;
}

export interface HeadcountRow {
  department: string;
  count: number;
}



export interface Paginated<T> {
  data: T[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

async function req<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers);
  if (init?.body && !headers.has("Content-Type")) headers.set("Content-Type", "application/json");

  if (typeof window !== "undefined") {
    const token = localStorage.getItem("token");
    if (token) headers.set("Authorization", `Bearer ${token}`);
  }

  const res = await fetch(`${BASE}${path}`, { ...init, headers });
  if (!res.ok) {
    if (res.status === 401 && typeof window !== "undefined" && window.location.pathname !== "/login" && window.location.pathname !== "/register") {
      localStorage.removeItem("token");
      window.location.href = "/login";
      return new Promise(() => {}); // Wait for redirect
    }

    let detail = "request failed";
    try {
      const body = await res.json();
      detail = body.detail || body.title || detail;
    } catch {
      detail = res.statusText || detail;
    }
    throw new Error(`${res.status} ${detail}`);
  }
  return res.json() as Promise<T>;
}

export const api = {
  listPersons: (q: Record<string, string | number | undefined> = {}) => {
    const sp = new URLSearchParams();
    for (const [k, v] of Object.entries(q)) if (v !== undefined && v !== "") sp.set(k, String(v));
    const qs = sp.toString();
    return req<Paginated<Person>>(`/api/v1/persons${qs ? `?${qs}` : ""}`);
  },
  getPerson: (id: string) => req<Person>(`/api/v1/persons/${id}`),
  createPerson: (body: unknown) => req<Person>(`/api/v1/persons`, { method: "POST", body: JSON.stringify(body) }),
  getEmergencyContact: (personId: string) => req<EmergencyContact>(`/api/v1/persons/${personId}/emergency-contact`),
  upsertEmergencyContact: (personId: string, body: unknown) => req<EmergencyContact>(`/api/v1/persons/${personId}/emergency-contact`, { method: "POST", body: JSON.stringify(body) }),
  listEmployment: (id: string) => req<EmploymentRecord[]>(`/api/v1/persons/${id}/employment`),
  getCurrentEmployment: (id: string) => req<EmploymentRecord>(`/api/v1/persons/${id}/employment/current`),
  listEvents: (id: string, q: Record<string, string | number | undefined> = {}) => {
    const sp = new URLSearchParams();
    for (const [k, v] of Object.entries(q)) if (v !== undefined && v !== "") sp.set(k, String(v));
    const qs = sp.toString();
    return req<Paginated<StatusChangeEvent>>(`/api/v1/persons/${id}/events${qs ? `?${qs}` : ""}`);
  },
  headcount: () => req<HeadcountRow[]>(`/api/v1/analytics/headcount`),
  health: () => req<{ status: string; database: string }>(`/health`),
};
