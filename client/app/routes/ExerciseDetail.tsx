import { useParams } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { getExercise } from "api/api";
import ActivityGrid from "~/components/ActivityGrid";

export default function ExerciseDetail() {
  const { id } = useParams();
  const exerciseId = parseInt(id || "0");

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["exercise", exerciseId],
    queryFn: () => getExercise(exerciseId),
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

  const pgTitle = `${data.name} - Genki`;

  return (
    <main className="w-full px-5 md:px-20 pt-10 min-h-screen">
      <title>{pgTitle}</title>
      <h1>{data.name}</h1>
      {data.category && (
        <p className="text-(--color-fg-tertiary) text-sm">{data.category.name}</p>
      )}
      {data.description && (
        <p className="mt-2 text-sm max-w-lg">{data.description}</p>
      )}

      {data.muscles && data.muscles.length > 0 && (
        <div className="mt-6">
          <h3>Muscles</h3>
          <div className="flex flex-wrap gap-2 mt-2">
            {data.muscles.map((m) => (
              <span
                key={m.id}
                className="px-3 py-1 rounded-full text-sm bg-(--color-bg-tertiary)"
              >
                {m.name_en || m.name}
              </span>
            ))}
          </div>
        </div>
      )}

      <div className="mt-6 flex flex-col sm:flex-row gap-10">
        <div>
          <h3>Stats</h3>
          <div className="mt-2">
            <div>
              <span className="header-font font-bold text-xl">
                {data.total_sets ?? 0}
              </span>{" "}
              Sets
            </div>
            <div>
              <span className="header-font font-bold text-xl">
                {data.total_reps ?? 0}
              </span>{" "}
              Reps
            </div>
            {data.total_volume_kg != null && data.total_volume_kg > 0 && (
              <div>
                <span className="header-font font-bold text-xl">
                  {data.total_volume_kg.toFixed(1)}
                </span>{" "}
                kg total volume
              </div>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
