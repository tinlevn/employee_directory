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

  if (error) return <p className="rounded bg-red-50 p-4 text-sm text-red-700">{error}</p>;
  if (!person) return <p className="text-sm text-zinc-500">Loading employee...</p>;

  return (
    <div className="space-y-6">
      <section className="rounded-lg border bg-white p-6 border-t-4 border-t-emerald-500 shadow-sm">
        <div className="flex items-center gap-4">
            {person.profile_photo_url ? (
                <img src={person.profile_photo_url} alt={`${person.first_name} ${person.last_name}`} className="w-16 h-16 rounded-full object-cover border border-zinc-200" width="64" height="64" />
            ) : (
                <div className="w-16 h-16 rounded-full bg-emerald-100 flex items-center justify-center text-emerald-800 font-semibold text-xl border border-emerald-200">
                    {person.first_name.charAt(0)}{person.last_name.charAt(0)}
                </div>
            )}
            <div>
                <h1 className="text-xl font-semibold text-zinc-900">{person.first_name} {person.last_name}</h1>
                <p className="mt-1 text-sm text-zinc-600">{employment?.job_title || "No position"} · {employment?.department || "No department"}</p>
            </div>
        </div>
        <div className="mt-6 grid gap-4 text-sm sm:grid-cols-2">
          <div><span className="text-zinc-500">Organization email</span><p>{person.org_email || "—"}</p></div>
          <div><span className="text-zinc-500">Phone</span><p>{person.phone_primary || "—"}</p></div>
          <div><span className="text-zinc-500">Location</span><p>{person.city || "—"}{person.country ? `, ${person.country}` : ""}</p></div>
          <div><span className="text-zinc-500">Hire date</span><p>{employment?.hire_date?.slice(0, 10) || "—"}</p></div>
        </div>
      </section>

      <section className="rounded-lg border bg-white p-6">
        <h2 className="font-medium">Current employment</h2>
        <div className="mt-4 grid gap-4 text-sm sm:grid-cols-3">
          <div><span className="text-zinc-500">Job level</span><p>{employment?.job_level || "—"}</p></div>
          <div><span className="text-zinc-500">Team</span><p>{employment?.team || "—"}</p></div>
          <div><span className="text-zinc-500">Work arrangement</span><p>{employment?.work_arrangement || "—"}</p></div>
        </div>
      </section>

      <section className="rounded-lg border bg-white p-6">
        <h2 className="font-medium">Emergency contact</h2>
        {contact ? (
          <div className="mt-3 grid gap-4 text-sm sm:grid-cols-3">
            <div><span className="text-zinc-500">Name</span><p>{contact.name}</p></div>
            <div><span className="text-zinc-500">Relationship</span><p>{contact.relationship || "—"}</p></div>
            <div><span className="text-zinc-500">Phone / email</span><p>{contact.phone || contact.email || "—"}</p></div>
          </div>
        ) : <p className="mt-3 text-sm text-zinc-500">No emergency contact recorded.</p>}
      </section>

      <section className="rounded-lg border bg-white p-6">
        <h2 className="font-medium">Event timeline</h2>
        {events.length ? (
          <div className="mt-4 divide-y divide-zinc-100 border rounded-md">
            {events.map((event) => (
              <div key={event.id} className="flex flex-wrap justify-between gap-2 px-4 py-3 text-sm hover:bg-emerald-50/50 transition-colors">
                <span className="font-medium text-zinc-700">{event.event_type}</span>
                <span className="text-zinc-500 bg-zinc-100 px-2 rounded-md">{event.effective_date.slice(0, 10)}</span>
              </div>
            ))}
          </div>
        ) : <p className="mt-3 text-sm text-zinc-500">No events recorded.</p>}
      </section>
    </div>
  );
}
