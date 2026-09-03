import { useEffect, useState } from "react";
import { api, type EmergencyContact, type EmploymentRecord, type Person, type StatusChangeEvent } from "../lib/api";

interface Props {
  id?: string;
}

export default function PersonDetail({ id }: Props) {
  const [person, setPerson] = useState<Person | null>(null);
  const [employment, setEmployment] = useState<EmploymentRecord | null>(null);
  const [contact, setContact] = useState<EmergencyContact | null>(null);
  const [events, setEvents] = useState<StatusChangeEvent[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const personId = id || new URLSearchParams(window.location.search).get("id") || "";
    if (!personId) {
      setError("No employee ID supplied.");
      return;
    }
    let cancelled = false;
    Promise.all([
      api.getPerson(personId),
      api.getCurrentEmployment(personId).catch(() => null),
      api.getEmergencyContact(personId).catch(() => null),
      api.listEvents(personId, { page_size: 20 }).then((result) => result.data).catch(() => []),
    ])
      .then(([loadedPerson, loadedEmployment, loadedContact, loadedEvents]) => {
        if (cancelled) return;
        setPerson(loadedPerson);
        setEmployment(loadedEmployment);
        setContact(loadedContact);
        setEvents(loadedEvents);
      })
      .catch((reason: Error) => {
        if (!cancelled) setError(reason.message);
      });
    return () => {
      cancelled = true;
    };
  }, [id]);

  if (error) return <p className="rounded-md bg-red-500/10 border border-red-500/20 p-4 text-sm text-red-600 dark:text-red-400">{error}</p>;
  if (!person) return <p className="text-sm text-[#5A6578] dark:text-slate-400">Loading employee...</p>;

  return (
    <div className="space-y-6">
      <section className="rounded-lg border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] p-6 border-t-4 border-t-[#41B06E] dark:border-t-[#1DCD9F] shadow-sm">
        <div className="flex items-center gap-4">
            {person.profile_photo_url ? (
                <img src={person.profile_photo_url} alt={`${person.first_name} ${person.last_name}`} className="w-16 h-16 rounded-full object-cover border border-[#E6DBC5] dark:border-[#2b303c]" width="64" height="64" />
            ) : (
                <div className="w-16 h-16 rounded-full bg-[#8DECB4]/30 dark:bg-[#1DCD9F]/20 flex items-center justify-center text-[#141E46] dark:text-[#1DCD9F] font-bold text-xl border border-[#41B06E]/30 dark:border-[#1DCD9F]/30">
                    {person.first_name.charAt(0)}{person.last_name.charAt(0)}
                </div>
            )}
            <div>
                <h1 className="text-xl font-semibold text-[#141E46] dark:text-white">{person.first_name} {person.last_name}</h1>
                <p className="mt-1 text-sm text-[#5A6578] dark:text-slate-400">{employment?.job_title || "No position"} · {employment?.department || "No department"}</p>
            </div>
        </div>
        <div className="mt-6 grid gap-4 text-sm sm:grid-cols-2">
          <div><span className="text-[#5A6578] dark:text-slate-400">Organization email</span><p className="font-medium text-[#141E46] dark:text-slate-100">{person.org_email || "—"}</p></div>
          <div><span className="text-[#5A6578] dark:text-slate-400">Phone</span><p className="font-medium text-[#141E46] dark:text-slate-100">{person.phone_primary || "—"}</p></div>
          <div><span className="text-[#5A6578] dark:text-slate-400">Location</span><p className="font-medium text-[#141E46] dark:text-slate-100">{person.city || "—"}{person.country ? `, ${person.country}` : ""}</p></div>
          <div><span className="text-[#5A6578] dark:text-slate-400">Hire date</span><p className="font-medium text-[#141E46] dark:text-slate-100">{employment?.hire_date?.slice(0, 10) || "—"}</p></div>
        </div>
      </section>

      <section className="rounded-lg border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] p-6 shadow-sm">
        <h2 className="font-semibold text-[#141E46] dark:text-slate-100 flex items-center gap-2">
          <span className="h-2 w-2 rounded-full bg-[#41B06E] dark:bg-[#1DCD9F]"></span>
          Current employment
        </h2>
        <div className="mt-4 grid gap-4 text-sm sm:grid-cols-3">
          <div><span className="text-[#5A6578] dark:text-slate-400">Job level</span><p className="font-medium text-[#141E46] dark:text-slate-100">{employment?.job_level || "—"}</p></div>
          <div><span className="text-[#5A6578] dark:text-slate-400">Team</span><p className="font-medium text-[#141E46] dark:text-slate-100">{employment?.team || "—"}</p></div>
          <div><span className="text-[#5A6578] dark:text-slate-400">Work arrangement</span><p className="font-medium text-[#141E46] dark:text-slate-100">{employment?.work_arrangement || "—"}</p></div>
        </div>
      </section>

      <section className="rounded-lg border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] p-6 shadow-sm">
        <h2 className="font-semibold text-[#141E46] dark:text-slate-100 flex items-center gap-2">
          <span className="h-2 w-2 rounded-full bg-[#41B06E] dark:bg-[#1DCD9F]"></span>
          Emergency contact
        </h2>
        {contact ? (
          <div className="mt-3 grid gap-4 text-sm sm:grid-cols-3">
            <div><span className="text-[#5A6578] dark:text-slate-400">Name</span><p className="font-medium text-[#141E46] dark:text-slate-100">{contact.name}</p></div>
            <div><span className="text-[#5A6578] dark:text-slate-400">Relationship</span><p className="font-medium text-[#141E46] dark:text-slate-100">{contact.relationship || "—"}</p></div>
            <div><span className="text-[#5A6578] dark:text-slate-400">Phone / email</span><p className="font-medium text-[#141E46] dark:text-slate-100">{contact.phone || contact.email || "—"}</p></div>
          </div>
        ) : <p className="mt-3 text-sm text-[#5A6578] dark:text-slate-400">No emergency contact recorded.</p>}
      </section>

      <section className="rounded-lg border border-[#E6DBC5] dark:border-[#2b303c] bg-white dark:bg-[#1c1f26] p-6 shadow-sm">
        <h2 className="font-semibold text-[#141E46] dark:text-slate-100 flex items-center gap-2">
          <span className="h-2 w-2 rounded-full bg-[#41B06E] dark:bg-[#1DCD9F]"></span>
          Event timeline
        </h2>
        {events.length ? (
          <div className="mt-4 divide-y divide-[#E6DBC5] dark:divide-[#2b303c] border border-[#E6DBC5] dark:border-[#2b303c] rounded-md overflow-hidden">
            {events.map((event) => (
              <div key={event.id} className="flex flex-wrap justify-between gap-2 px-4 py-3 text-sm hover:bg-[#8DECB4]/20 dark:hover:bg-[#1DCD9F]/10 transition-colors">
                <span className="font-medium text-[#141E46] dark:text-slate-200">{event.event_type}</span>
                <span className="text-[#5A6578] dark:text-slate-400 bg-[#F8EFE0] dark:bg-[#252a34] px-2 py-0.5 rounded text-xs font-mono">{event.effective_date.slice(0, 10)}</span>
              </div>
            ))}
          </div>
        ) : <p className="mt-3 text-sm text-[#5A6578] dark:text-slate-400">No events recorded.</p>}
      </section>
    </div>
  );
}
