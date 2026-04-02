import { useParams } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { getWorkout } from "api/api";
import { Link } from "react-router";

export default function WorkoutDetail() {
  const { id } = useParams();
  const workoutId = parseInt(id || "0");

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["workout", workoutId],
    queryFn: () => getWorkout(workoutId),
  });

  if (isPending) {
    return (
      <main className="w-full px-5 md:px-20 pt-10">
        <p>Loading...</p>
      </main>
    );
  }
  if (isError) {
    return (
      <main className="w-full px-5 md:px-20 pt-10">
        <p className="error">Error: {error.message}</p>
      </main>
    );
  }

  const { workout, sets } = data;
  const title = workout.title || "Workout";
  const pgTitle = `${title} - Genki`;
  const startDate = new Date(workout.started_at);

  // Group sets by exercise
  const byExercise = new Map<number, typeof sets>();
  for (const s of sets) {
    const eid = s.exercise_id;
    if (!byExercise.has(eid)) byExercise.set(eid, []);
    byExercise.get(eid)!.push(s);
  }

  return (
    <main className="w-full px-5 md:px-20 pt-10 min-h-screen">
      <title>{pgTitle}</title>
      <h1>{title}</h1>
      <p className="text-(--color-fg-tertiary) text-sm">
        {startDate.toLocaleDateString()} at {startDate.toLocaleTimeString()}
        {workout.duration_minutes && ` · ${workout.duration_minutes} min`}
        {` · ${workout.source}`}
      </p>
      {workout.notes && <p className="mt-2 text-sm">{workout.notes}</p>}

      <div className="mt-8 flex flex-col gap-6">
        {Array.from(byExercise.entries()).map(([eid, exSets]) => (
          <div key={eid}>
            <h3>
              <Link
                to={`/exercise/${eid}`}
                className="hover:text-(--color-fg-secondary)"
              >
                {exSets[0].exercise?.name || `Exercise #${eid}`}
              </Link>
            </h3>
            <table className="mt-2 text-sm">
              <thead>
                <tr className="text-(--color-fg-tertiary)">
                  <th className="pr-4 text-left">Set</th>
                  <th className="pr-4 text-left">Reps</th>
                  <th className="pr-4 text-left">Weight</th>
                  <th className="pr-4 text-left">RPE</th>
                </tr>
              </thead>
              <tbody>
                {exSets.map((s) => (
                  <tr key={s.id}>
                    <td className="pr-4">{s.set_number}</td>
                    <td className="pr-4">{s.reps ?? "-"}</td>
                    <td className="pr-4">
                      {s.weight_kg ? `${s.weight_kg} kg` : "-"}
                    </td>
                    <td className="pr-4">{s.rpe ?? "-"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ))}
      </div>
    </main>
  );
}
