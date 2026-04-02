import { useQuery } from "@tanstack/react-query";
import { getWorkouts, type getItemsArgs } from "api/api";
import { Link } from "react-router";
import { timeSince } from "~/utils/utils";
import { useState } from "react";

export default function WorkoutList() {
  const [page, setPage] = useState(1);

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["workouts-list", { limit: 50, period: "all_time", page }],
    queryFn: ({ queryKey }) => getWorkouts(queryKey[1] as getItemsArgs),
  });

  const pgTitle = "Workouts - Genki";

  return (
    <div className="w-full min-h-screen">
      <title>{pgTitle}</title>
      <div className="w-19/20 sm:w-17/20 mx-auto pt-6 sm:pt-12">
        <h1>Workouts</h1>
        <div className="mt-6">
          {isPending ? (
            <p>Loading...</p>
          ) : isError ? (
            <p className="error">Error: {error.message}</p>
          ) : (
            <>
              <table className="text-sm w-full">
                <thead>
                  <tr className="text-(--color-fg-tertiary) border-b border-(--color-bg-tertiary)">
                    <th className="text-left py-2">Date</th>
                    <th className="text-left py-2">Title</th>
                    <th className="text-left py-2">Duration</th>
                    <th className="text-left py-2">Source</th>
                  </tr>
                </thead>
                <tbody>
                  {data.items.map((w) => (
                    <tr
                      key={w.id}
                      className="border-b border-(--color-bg-secondary) hover:bg-(--color-bg-secondary)"
                    >
                      <td className="py-2 pr-4 whitespace-nowrap text-(--color-fg-tertiary)">
                        {timeSince(new Date(w.started_at))}
                      </td>
                      <td className="py-2">
                        <Link
                          to={`/workout/${w.id}`}
                          className="hover:text-(--color-fg-secondary)"
                        >
                          {w.title || "Workout"}
                        </Link>
                      </td>
                      <td className="py-2 text-(--color-fg-tertiary)">
                        {w.duration_minutes ? `${w.duration_minutes} min` : "-"}
                      </td>
                      <td className="py-2 text-(--color-fg-tertiary) text-xs">
                        {w.source}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              <div className="flex gap-4 mt-6">
                <button
                  disabled={page <= 1}
                  onClick={() => setPage(page - 1)}
                  className="px-3 py-1 rounded border border-gray-400 disabled:opacity-30"
                >
                  Prev
                </button>
                <button
                  disabled={!data.has_next_page}
                  onClick={() => setPage(page + 1)}
                  className="px-3 py-1 rounded border border-gray-400 disabled:opacity-30"
                >
                  Next
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
