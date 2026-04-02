import { useLocation } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { getTopExercises, type getItemsArgs } from "api/api";
import TopItemList from "~/components/TopItemList";
import PeriodSelector from "~/components/PeriodSelector";
import { useState } from "react";

export default function ExerciseChart() {
  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const [period, setPeriod] = useState(params.get("period") || "week");
  const [page, setPage] = useState(1);

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["top-exercises-chart", { limit: 50, period, page }],
    queryFn: ({ queryKey }) => getTopExercises(queryKey[1] as getItemsArgs),
  });

  const pgTitle = "Top Exercises - Genki";

  return (
    <div className="w-full min-h-screen">
      <title>{pgTitle}</title>
      <div className="w-19/20 sm:w-17/20 mx-auto pt-6 sm:pt-12">
        <h1>Top Exercises</h1>
        <PeriodSelector current={period} setter={setPeriod} disableCache />
        <div className="mt-10">
          {isPending ? (
            <p>Loading...</p>
          ) : isError ? (
            <p className="error">Error: {error.message}</p>
          ) : (
            <>
              <TopItemList type="exercise" data={data} ranked separators />
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
