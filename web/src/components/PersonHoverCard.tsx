import { useEffect, useState, type ReactNode } from "react";
import { api, type EmploymentRecord, type Person } from "../lib/api";

interface Details {
  person: Person;
  employment: EmploymentRecord | null;
}

const cache = new Map<string, Details>();
const CARD_WIDTH = 340;
const EST_HEIGHT = 380;

interface Props {
  personId: string;
  fallback: Person;
  anchorRect: DOMRect;
  onEnter: () => void;
  onLeave: () => void;
}

function initials(p: Person) {
  return `${p.first_name?.[0] ?? ""}${p.last_name?.[0] ?? ""}`.toUpperCase();
}

function Row({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="flex justify-between gap-3">
      <span className="shrink-0 text-zinc-500">{label}</span>
      <span className="min-w-0 truncate text-right text-zinc-800">{children}</span>
    </div>
  );
}

export default function PersonHoverCard({ personId, fallback, anchorRect, onEnter, onLeave }: Props) {
  const [details, setDetails] = useState<Details | null>(cache.get(personId) ?? null);

  useEffect(() => {
    const cached = cache.get(personId);
    if (cached) {
      setDetails(cached);
      return;
    }
    let cancelled = false;
    setDetails(null);
    Promise.all([api.getPerson(personId), api.getCurrentEmployment(personId).catch(() => null)])
      .then(([person, employment]) => {
        const entry: Details = { person, employment };
        cache.set(personId, entry);
        if (!cancelled) setDetails(entry);
      })
      .catch(() => {
        if (!cancelled) setDetails(null);
      });
    return () => {
      cancelled = true;
    };
  }, [personId]);

  const person = details?.person ?? fallback;
  const employment = details?.employment ?? null;
  const loading = !details;

  const vw = window.innerWidth;
  const vh = window.innerHeight;
  const left = Math.min(Math.max(anchorRect.left, 8), vw - CARD_WIDTH - 8);
  const below = vh - anchorRect.bottom >= EST_HEIGHT;

  return (
    <div
      onMouseEnter={onEnter}
      onMouseLeave={onLeave}
      style={{
        position: "fixed",
        left,
        width: CARD_WIDTH,
        ...(below ? { top: anchorRect.bottom + 8 } : { bottom: vh - anchorRect.top + 8 }),
      }}
      className="z-50 rounded-lg border bg-white p-4 shadow-xl"
    >
      <div className="flex items-start gap-3">
        {person.profile_photo_url ? (
          <img src={person.profile_photo_url} alt="" className="h-12 w-12 shrink-0 rounded-full object-cover" />
        ) : (
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-zinc-900 text-sm font-semibold text-white">
            {initials(person)}
          </div>
        )}
        <div className="min-w-0 flex-1">
          <p className="truncate font-semibold text-zinc-900">
            {person.first_name} {person.last_name}
            {loading && <span className="ml-2 text-xs font-normal text-zinc-400">loading…</span>}
          </p>
          {person.preferred_name && <p className="text-xs text-zinc-500">&ldquo;{person.preferred_name}&rdquo;</p>}
          <p className="mt-0.5 truncate text-sm text-zinc-600">
            {employment?.job_title || person.current_job_title || "—"}
          </p>
          <div className="mt-1.5 flex flex-wrap gap-1">
            {(employment?.department || person.current_department) && (
              <span className="rounded-full bg-zinc-900 px-2.5 py-0.5 text-xs font-medium text-white">
                {employment?.department || person.current_department}
              </span>
            )}
            {employment?.employment_status && (
              <span className="rounded-full bg-zinc-100 px-2 py-0.5 text-xs text-zinc-600">
                {employment.employment_status}
              </span>
            )}
            {!person.is_active && (
              <span className="rounded-full bg-red-50 px-2 py-0.5 text-xs text-red-700">inactive</span>
            )}
          </div>
        </div>
      </div>

      <div className="mt-3 space-y-1.5 border-t pt-3 text-sm">
        <Row label="Email">{person.org_email || person.personal_email || "—"}</Row>
        <Row label="Phone">{person.phone_primary || "—"}</Row>
        <Row label="Location">{[person.city, person.country].filter(Boolean).join(", ") || "—"}</Row>
        <Row label="Office">{person.current_office_location || "—"}</Row>
        <Row label="Team">{employment?.team || person.current_team || "—"}</Row>
        <Row label="Hired">{(employment?.hire_date || person.current_hire_date || "").slice(0, 10) || "—"}</Row>
        {employment?.job_level && <Row label="Level">{employment.job_level}</Row>}
      </div>

      {person.tags.length > 0 && (
        <div className="mt-3 flex flex-wrap gap-1 border-t pt-3">
          {person.tags.map((t) => (
            <span key={t} className="rounded-full bg-zinc-100 px-2 py-0.5 text-xs text-zinc-600">
              {t}
            </span>
          ))}
        </div>
      )}

      <a
        href={`/person?id=${encodeURIComponent(personId)}`}
        className="mt-3 block rounded-md bg-zinc-900 px-3 py-2 text-center text-sm font-medium text-white hover:bg-zinc-800"
      >
        View full profile →
      </a>
    </div>
  );
}
